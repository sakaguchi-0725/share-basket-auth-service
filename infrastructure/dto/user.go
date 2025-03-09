package dto

import (
	"share-basket-auth-service/domain/model"
	"time"
)

type User struct {
	ID         string    `gorm:"primaryKey"`
	CognitoUID string    `gorm:"not null"`
	Email      string    `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

func NewUserDto(user *model.User) User {
	return User{
		ID:         user.ID.String(),
		CognitoUID: user.CognitoUID,
		Email:      user.Email,
	}
}

func (u *User) ToModel() model.User {
	id, _ := model.NewUserID(u.ID)

	return model.RecreateUser(
		id, u.CognitoUID, u.Email,
	)
}
