package model

import "errors"

type User struct {
	ID         UserID
	CognitoUID string
	Email      string
}

func NewUser(id UserID, cognitoUID, email string) (User, error) {
	if email == "" {
		return User{}, errors.New("email is required")
	}

	if cognitoUID == "" {
		return User{}, errors.New("cognitoUID is required")
	}

	return User{
		ID:         id,
		CognitoUID: cognitoUID,
		Email:      email,
	}, nil
}

func RecreateUser(id UserID, cognitoUID, email string) User {
	return User{
		ID:         id,
		CognitoUID: cognitoUID,
		Email:      email,
	}
}
