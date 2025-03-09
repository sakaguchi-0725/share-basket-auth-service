//go:generate mockgen -source=$GOFILE -destination=../tests/mock/$GOPACKAGE/mock_$GOFILE -package=mock
package usecase

import (
	"context"
	"errors"
	"share-basket-auth-service/core/apperr"
	"share-basket-auth-service/domain/model"
	"share-basket-auth-service/domain/repository"
	"share-basket-auth-service/domain/service"
)

type SignUpUseCase interface {
	Execute(ctx context.Context, email, password string) error
}

type signUpUseCase struct {
	authenticator repository.Authenticator
	userRepo      repository.User
	userService   service.User
}

func (s *signUpUseCase) Execute(ctx context.Context, email string, password string) error {
	available, err := s.userService.IsEmailAvailable(email)
	if err != nil {
		return err
	}

	if !available {
		return apperr.NewApplicationError(apperr.ErrBadRequest, apperr.ErrInvalidData.Error(), apperr.ErrInvalidData)
	}

	cognitoUID, err := s.authenticator.SignUp(ctx, email, password)
	if err != nil {
		if errors.Is(err, apperr.ErrDuplicatedKey) || errors.Is(err, apperr.ErrInvalidData) {
			return apperr.NewApplicationError(apperr.ErrBadRequest, apperr.ErrInvalidData.Error(), err)
		}
		return err
	}

	user, err := model.NewUser(model.GenerateUserID(), cognitoUID, email)
	if err != nil {
		return apperr.NewApplicationError(apperr.ErrBadRequest, err.Error(), err)
	}

	if err := s.userRepo.Create(ctx, &user); err != nil {
		return err
	}

	return nil
}

func NewSignUpUseCase(
	authenticator repository.Authenticator,
	userRepo repository.User,
	userService service.User,
) SignUpUseCase {
	return &signUpUseCase{authenticator, userRepo, userService}
}
