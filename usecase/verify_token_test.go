package usecase_test

import (
	"context"
	"share-basket-auth-service/core/apperr"
	"share-basket-auth-service/domain/model"
	. "share-basket-auth-service/tests/mock/repository"
	"share-basket-auth-service/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestVerifyToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	authenticator := NewMockAuthenticator(ctrl)
	userRepo := NewMockUser(ctrl)
	usecase := usecase.NewVerifyTokenUseCase(authenticator, userRepo)

	defer ctrl.Finish()

	tests := map[string]struct {
		setupMock   func()
		token       string
		expected    string
		expectedErr error
	}{
		"正常系": {
			setupMock: func() {
				authenticator.EXPECT().VerifyToken(gomock.Any(), "dummy-token").Return("test@example.com", nil)
				userRepo.EXPECT().GetByEmail(gomock.Any(), "test@example.com").Return(model.User{ID: model.UserID("dummy-user-id")}, nil)
			},
			token:       "dummy-token",
			expected:    "dummy-user-id",
			expectedErr: nil,
		},
		"トークンが無効": {
			setupMock: func() {
				authenticator.EXPECT().VerifyToken(gomock.Any(), "invalid-token").Return("", apperr.ErrInvalidToken)
			},
			token:       "invalid-token",
			expected:    "",
			expectedErr: apperr.ErrInvalidToken,
		},
		"ユーザーが存在しない": {
			setupMock: func() {
				authenticator.EXPECT().VerifyToken(gomock.Any(), "dummy-token").Return("test@example.com", nil)
				userRepo.EXPECT().GetByEmail(gomock.Any(), "test@example.com").Return(model.User{}, apperr.ErrDataNotFound)
			},
			token:       "dummy-token",
			expected:    "",
			expectedErr: apperr.ErrDataNotFound,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.setupMock()
			actual, err := usecase.Execute(context.TODO(), tt.token)

			assert.Equal(t, tt.expected, actual)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
