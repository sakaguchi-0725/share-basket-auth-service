//go:generate mockgen -source=$GOFILE -destination=../../tests/mock/$GOPACKAGE/mock_$GOFILE -package=mock
package service

import (
	"context"
	"share-basket-auth-service/core/apperr"
	"share-basket-auth-service/domain/repository"
)

type User interface {
	IsEmailAvailable(email string) (bool, error)
}

type userService struct {
	userRepo repository.User
}

func (u *userService) IsEmailAvailable(email string) (bool, error) {
	_, err := u.userRepo.GetByEmail(context.Background(), email)
	if err != nil {
		if err == apperr.ErrDataNotFound {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func NewUserService(userRepo repository.User) User {
	return &userService{userRepo}
}
