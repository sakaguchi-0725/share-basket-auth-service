package usecase_test

import (
	"context"
	"share-basket-auth-service/core/apperr"
	. "share-basket-auth-service/tests/mock/repository"
	"share-basket-auth-service/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestLoginUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	authenticator := NewMockAuthenticator(ctrl)
	usecase := usecase.NewLoginUseCase(authenticator)

	defer ctrl.Finish()

	tests := map[string]struct {
		setupMock   func()
		email       string
		password    string
		expected    string
		expectedErr error
	}{
		"正常系": {
			setupMock: func() {
				authenticator.EXPECT().Login(gomock.Any(), "test@example.com", "password").Return("dummy-token", nil)
			},
			email:       "test@example.com",
			password:    "password",
			expected:    "dummy-token",
			expectedErr: nil,
		},
		"認証に失敗した場合": {
			setupMock: func() {
				authenticator.EXPECT().Login(gomock.Any(), "test@example.com", "password").Return("", apperr.ErrUnauthenticated)
			},
			email:       "test@example.com",
			password:    "password",
			expected:    "",
			expectedErr: apperr.ErrUnauthenticated,
		},
		"入力値が不正な場合": {
			setupMock: func() {
				authenticator.EXPECT().Login(gomock.Any(), "test@example.com", "password").Return("", apperr.ErrInvalidData)
			},
			email:       "test@example.com",
			password:    "password",
			expected:    "",
			expectedErr: apperr.ErrInvalidData,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.setupMock()

			actual, err := usecase.Execute(context.Background(), tt.email, tt.password)

			assert.Equal(t, tt.expected, actual)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, tt.expectedErr, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
