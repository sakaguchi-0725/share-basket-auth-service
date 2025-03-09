//go:generate mockgen -source=$GOFILE -destination=../tests/mock/$GOPACKAGE/mock_$GOFILE -package=mock
package usecase

import (
	"context"
	"errors"
	"share-basket-auth-service/core/apperr"
	"share-basket-auth-service/domain/repository"
)

type SignUpConfirmUseCase interface {
	Execute(ctx context.Context, email, confirmationCode string) error
}

type signUpConfirmUseCase struct {
	authenticator repository.Authenticator
}

func (s *signUpConfirmUseCase) Execute(ctx context.Context, email string, confirmationCode string) error {
	err := s.authenticator.SignUpConfirm(ctx, email, confirmationCode)
	if err != nil {
		if errors.Is(err, apperr.ErrInvalidData) || errors.Is(err, apperr.ErrExpiredCodeException) {
			return apperr.NewApplicationError(apperr.ErrBadRequest, err.Error(), err)
		}
		return err
	}

	return nil
}

func NewSignUpConfirmUseCase(authenticator repository.Authenticator) SignUpConfirmUseCase {
	return &signUpConfirmUseCase{authenticator}
}
