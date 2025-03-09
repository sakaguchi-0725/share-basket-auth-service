//go:generate mockgen -source=$GOFILE -destination=../tests/mock/$GOPACKAGE/mock_$GOFILE -package=mock
package usecase

import (
	"context"
	"errors"
	"share-basket-auth-service/core/apperr"
	"share-basket-auth-service/domain/repository"
)

type VerifyTokenUseCase interface {
	Execute(ctx context.Context, token string) (string, error)
}

type verifyTokenUseCase struct {
	authenticator repository.Authenticator
	userRepo      repository.User
}

func (v *verifyTokenUseCase) Execute(ctx context.Context, token string) (string, error) {
	email, err := v.authenticator.VerifyToken(ctx, token)
	if err != nil {
		if errors.Is(err, apperr.ErrInvalidToken) || errors.Is(err, apperr.ErrTokenExpired) {
			return "", apperr.NewApplicationError(apperr.ErrUnauthorized, err.Error(), err)
		}

		return "", err
	}

	user, err := v.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, apperr.ErrDataNotFound) {
			return "", apperr.NewApplicationError(apperr.ErrUnauthorized, err.Error(), err)
		}

		return "", err
	}

	return user.ID.String(), nil
}

func NewVerifyTokenUseCase(authenticator repository.Authenticator, userRepo repository.User) VerifyTokenUseCase {
	return &verifyTokenUseCase{authenticator, userRepo}
}
