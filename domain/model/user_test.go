package model_test

import (
	"errors"
	"share-basket-auth-service/domain/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	t.Run("NewUser", func(t *testing.T) {
		userID := model.GenerateUserID()

		tests := map[string]struct {
			id          model.UserID
			cognitoUID  string
			email       string
			expected    model.User
			expectedErr error
		}{
			"success": {
				id:         userID,
				cognitoUID: "cognito-uid",
				email:      "test@example.com",
				expected: model.User{
					ID:         userID,
					CognitoUID: "cognito-uid",
					Email:      "test@example.com",
				},
				expectedErr: nil,
			},
			"missing email": {
				id:          userID,
				cognitoUID:  "cognito-uid",
				email:       "",
				expected:    model.User{},
				expectedErr: errors.New("email is required"),
			},
			"missing cognitoUID": {
				id:          userID,
				cognitoUID:  "",
				email:       "test@example.com",
				expected:    model.User{},
				expectedErr: errors.New("cognitoUID is required"),
			},
		}

		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				actual, err := model.NewUser(tt.id, tt.cognitoUID, tt.email)

				assert.Equal(t, tt.expected, actual)
				if tt.expectedErr != nil {
					assert.EqualError(t, err, tt.expectedErr.Error())
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}
