package model_test

import (
	"errors"
	"share-basket-auth-service/domain/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserID(t *testing.T) {
	t.Run("NewUserID", func(t *testing.T) {
		tests := map[string]struct {
			input       string
			expected    model.UserID
			expectedErr error
		}{
			"valid": {
				input:       "123e4567-e89b-12d3-a456-426614174000",
				expected:    model.UserID("123e4567-e89b-12d3-a456-426614174000"),
				expectedErr: nil,
			},
			"invalid": {
				input:       "invalid-uuid",
				expected:    model.UserID(""),
				expectedErr: errors.New("invalid user ID: invalid UUID length: 12"),
			},
		}

		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				actual, err := model.NewUserID(tt.input)

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
