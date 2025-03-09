package persistence

import (
	"context"
	"errors"
	"log"
	"share-basket-auth-service/core/apperr"
	"share-basket-auth-service/domain/model"
	"share-basket-auth-service/domain/repository"
	"share-basket-auth-service/infrastructure/dto"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func (ur *userRepository) Create(ctx context.Context, user *model.User) error {
	userDto := dto.NewUserDto(user)

	err := ur.db.Create(&userDto).Error
	if err != nil {
		return err
	}

	return nil
}

func (ur *userRepository) GetByEmail(ctx context.Context, email string) (model.User, error) {
	var userDto dto.User

	err := ur.db.Where("email = ?", email).First(&userDto).Error
	if err != nil {
		log.Println("Errorr fetching user by email:", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, apperr.ErrDataNotFound
		}
		return model.User{}, err
	}

	return userDto.ToModel(), nil
}

func NewUserRepository(db *gorm.DB) repository.User {
	return &userRepository{db}
}
