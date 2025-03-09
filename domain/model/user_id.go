package model

import (
	"fmt"

	"github.com/google/uuid"
)

type UserID string

func GenerateUserID() UserID {
	return UserID(uuid.NewString())
}

func NewUserID(s string) (UserID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return "", fmt.Errorf("invalid user ID: %w", err)
	}

	return UserID(id.String()), nil
}

func (u UserID) String() string {
	return string(u)
}
