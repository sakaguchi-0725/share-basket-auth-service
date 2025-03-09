//go:generate mockgen -source=$GOFILE -destination=../tests/mock/$GOPACKAGE/mock_$GOFILE -package=mock
package usecase

import (
	"context"
	"errors"
	"share-basket-auth-service/core/apperr"
	"share-basket-auth-service/domain/repository"
)

type LoginUseCase interface {
	Execute(ctx context.Context, email, password string) (string, error)
}

type loginUseCase struct {
	authenticator repository.Authenticator
}

func (l *loginUseCase) Execute(ctx context.Context, email string, password string) (string, error) {
	token, err := l.authenticator.Login(ctx, email, password)
	if err != nil {
		if errors.Is(err, apperr.ErrUnauthenticated) || errors.Is(err, apperr.ErrInvalidData) {
			return "", apperr.NewApplicationError(apperr.ErrUnauthorized, "メールアドレスまたはパスワードが間違っています。", err)
		}
		return "", err
	}

	return token, nil
}

func NewLoginUseCase(authenticator repository.Authenticator) LoginUseCase {
	return &loginUseCase{authenticator}
}
