package persistence

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"share-basket-auth-service/core/apperr"
	"share-basket-auth-service/core/config"
	"share-basket-auth-service/core/util"
	"share-basket-auth-service/domain/repository"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/golang-jwt/jwt/v5"
)

var (
	usernameExistsException   = &types.UsernameExistsException{}
	invalidPasswordException  = &types.InvalidPasswordException{}
	invalidParameterException = &types.InvalidParameterException{}
	expiredCodeException      = &types.ExpiredCodeException{}
	codeMismatchException     = &types.CodeMismatchException{}
	notAuthorizedException    = &types.NotAuthorizedException{}
)

type cognito struct {
	client       *cognitoidentityprovider.Client
	clientID     string
	clientSecret string
	userPoolID   string
	jwksURL      string
	keys         map[string]*rsa.PublicKey
	mu           sync.Mutex
}

func (c *cognito) Login(ctx context.Context, email string, password string) (string, error) {
	input := &cognitoidentityprovider.InitiateAuthInput{
		ClientId: aws.String(c.clientID),
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		AuthParameters: map[string]string{
			"USERNAME": email,
			"PASSWORD": password,
		},
	}

	result, err := c.client.InitiateAuth(ctx, input)
	if err != nil {
		if errors.Is(err, notAuthorizedException) {
			return "", apperr.ErrUnauthenticated
		}

		if errors.Is(err, invalidParameterException) {
			return "", apperr.ErrInvalidData
		}

		return "", err
	}

	accessToken := result.AuthenticationResult.AccessToken
	if accessToken == nil {
		return "", apperr.ErrUnauthenticated
	}

	return util.Derefer(accessToken), nil
}

func (c *cognito) SignUp(ctx context.Context, email string, password string) (string, error) {
	result, err := c.client.SignUp(ctx, &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(c.clientID),
		Username: aws.String(email),
		Password: aws.String(password),
	})
	if err != nil {
		if errors.Is(err, usernameExistsException) {
			return "", apperr.ErrDuplicatedKey
		}

		if errors.Is(err, invalidPasswordException) || errors.Is(err, invalidParameterException) {
			return "", apperr.ErrInvalidData
		}
		return "", err
	}

	cognitoUID := result.UserSub
	if cognitoUID == nil {
		return "", apperr.ErrUnauthenticated
	}

	return util.Derefer(cognitoUID), nil
}

func (c *cognito) SignUpConfirm(ctx context.Context, email string, confirmationCode string) error {
	_, err := c.client.ConfirmSignUp(ctx, &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(c.clientID),
		Username:         aws.String(email),
		ConfirmationCode: aws.String(confirmationCode),
	})
	if err != nil {
		if errors.Is(err, invalidParameterException) || errors.Is(err, codeMismatchException) {
			return apperr.ErrInvalidData
		}

		if errors.Is(err, expiredCodeException) {
			return apperr.ErrExpiredCodeException
		}

		return err
	}

	return nil
}

func (c *cognito) VerifyToken(ctx context.Context, token string) (string, error) {
	parsedToken, err := jwt.Parse(token, c.keyFunc)
	if err != nil {
		return "", apperr.ErrInvalidToken
	}

	if !parsedToken.Valid {
		return "", apperr.ErrInvalidToken
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", apperr.ErrInvalidToken
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return "", apperr.ErrInvalidToken
	}

	if int64(exp) < time.Now().Unix() {
		return "", apperr.ErrTokenExpired
	}

	email, ok := claims["username"].(string)
	if !ok {
		return "", apperr.ErrInvalidToken
	}

	return email, nil
}

// JWTの`kid`に対応する公開鍵を取得する
func (c *cognito) keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("missing kid in token header")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if key, exists := c.keys[kid]; exists {
		return key, nil
	}

	key, err := c.fetchJWKS(kid)
	if err != nil {
		return nil, err
	}

	c.keys[kid] = key
	return key, nil
}

// CognitoのJWKSエンドポイントから公開鍵を取得する
func (c *cognito) fetchJWKS(kid string) (*rsa.PublicKey, error) {
	resp, err := http.Get(c.jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	var jwks struct {
		Keys []struct {
			Kty string `json:"kty"`
			Kid string `json:"kid"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("failed to decode JWKS: %w", err)
	}

	for _, key := range jwks.Keys {
		if key.Kid == kid {
			nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
			if err != nil {
				return nil, fmt.Errorf("failed to decode N: %w", err)
			}

			eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
			if err != nil {
				return nil, fmt.Errorf("failed to decode E: %w", err)
			}

			e := 0
			for _, b := range eBytes {
				e = e<<8 + int(b)
			}

			rsaKey := &rsa.PublicKey{
				N: new(big.Int).SetBytes(nBytes),
				E: e,
			}

			return rsaKey, nil
		}
	}

	return nil, errors.New("public key not found in JWKS")
}

func NewCognito(ctx context.Context, conf config.AWSConfig) (repository.Authenticator, error) {
	cfg, err := awsConfig.LoadDefaultConfig(ctx, awsConfig.WithRegion(conf.Region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &cognito{
		client:       cognitoidentityprovider.NewFromConfig(cfg),
		clientID:     conf.ClientID,
		clientSecret: conf.ClientSecret,
		userPoolID:   conf.UserPoolID,
		jwksURL:      fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", conf.Region, conf.UserPoolID),
	}, nil
}
