package persistence_test

import (
	"context"
	"share-basket-auth-service/core/apperr"
	"share-basket-auth-service/domain/model"
	"share-basket-auth-service/infrastructure/dto"
	"share-basket-auth-service/infrastructure/persistence"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserPersistence(t *testing.T) {
	repo := persistence.NewUserRepository(testDB)

	t.Run("Create", func(t *testing.T) {
		defer clearTestData()

		tests := map[string]struct {
			input       *model.User
			expectedErr error
		}{
			"正常系": {
				input: &model.User{
					ID:         model.GenerateUserID(),
					CognitoUID: "test-cognito-uid",
					Email:      "test@example.com",
				},
				expectedErr: nil,
			},
		}

		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				err := repo.Create(context.Background(), tt.input)
				if err != tt.expectedErr {
					assert.Error(t, err)
					assert.EqualError(t, err, tt.expectedErr.Error())
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("getByEmail", func(t *testing.T) {
		defer clearTestData()
		userID := model.GenerateUserID()

		tests := map[string]struct {
			email       string
			expected    model.User
			expectedErr error
		}{
			"正常系": {
				email: "test@example.com",
				expected: model.User{
					ID:         userID,
					CognitoUID: "test-cognito-uid",
					Email:      "test@example.com",
				},
				expectedErr: nil,
			},
			"ユーザーが存在しない場合": {
				email:       "nonexistent@example.com",
				expected:    model.User{},
				expectedErr: apperr.ErrDataNotFound,
			},
		}

		t.Run("mocking", func(t *testing.T) {
			err := testDB.Create(&dto.User{
				ID:         userID.String(),
				CognitoUID: "test-cognito-uid",
				Email:      "test@example.com",
			}).Error
			require.NoError(t, err)
		})

		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				actual, err := repo.GetByEmail(context.Background(), tt.email)

				assert.Equal(t, tt.expected, actual)
				if tt.expectedErr != nil {
					assert.Error(t, err)
					assert.EqualError(t, err, tt.expectedErr.Error())
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}
