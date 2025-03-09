package service_test

import (
	"errors"
	"share-basket-auth-service/core/apperr"
	"share-basket-auth-service/domain/model"
	"share-basket-auth-service/domain/service"
	. "share-basket-auth-service/tests/mock/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserService(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepo := NewMockUser(ctrl)
	userService := service.NewUserService(userRepo)

	tests := map[string]struct {
		setupMock   func()
		email       string
		expected    bool
		expectedErr error
	}{
		"正常系: メールアドレスが存在しない場合": {
			setupMock: func() {
				userRepo.EXPECT().GetByEmail(gomock.Any(), "test@example.com").Return(model.User{}, apperr.ErrDataNotFound)
			},
			email:       "test@example.com",
			expected:    true,
			expectedErr: nil,
		},
		"正常系: メールアドレスが存在する場合": {
			setupMock: func() {
				userRepo.EXPECT().GetByEmail(gomock.Any(), "test@example.com").Return(model.User{
					ID:    "user-id",
					Email: "test@example.com",
				}, nil)
			},
			email:       "test@example.com",
			expected:    false,
			expectedErr: nil,
		},
		"データベースエラーが発生した場合": {
			setupMock: func() {
				userRepo.EXPECT().GetByEmail(gomock.Any(), "test@example.com").Return(model.User{}, errors.New("database error"))
			},
			email:       "test@example.com",
			expected:    false,
			expectedErr: errors.New("database error"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.setupMock()
			actual, err := userService.IsEmailAvailable("test@example.com")

			assert.Equal(t, tt.expected, actual)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
