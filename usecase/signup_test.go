package usecase_test

import (
	"context"
	"errors"
	"share-basket-auth-service/core/apperr"
	"share-basket-auth-service/domain/model"
	repo "share-basket-auth-service/tests/mock/repository"
	mock "share-basket-auth-service/tests/mock/service"
	"share-basket-auth-service/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSignUpUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	authenticator := repo.NewMockAuthenticator(ctrl)
	userRepo := repo.NewMockUser(ctrl)
	userService := mock.NewMockUser(ctrl)
	usecase := usecase.NewSignUpUseCase(authenticator, userRepo, userService)

	defer ctrl.Finish()

	tests := map[string]struct {
		setupMock   func()
		email       string
		password    string
		expectedErr error
	}{
		"正常系": {
			setupMock: func() {
				userService.EXPECT().IsEmailAvailable(gomock.Any()).Return(true, nil)
				authenticator.EXPECT().SignUp(gomock.Any(), "test@example.com", "password").Return("user-id", nil)
				userRepo.EXPECT().Create(gomock.Any(), gomock.AssignableToTypeOf(&model.User{})).Return(nil)
			},
			email:       "test@example.com",
			password:    "password",
			expectedErr: nil,
		},
		"メールアドレスが既に存在する場合": {
			setupMock: func() {
				userService.EXPECT().IsEmailAvailable(gomock.Any()).Return(false, nil)
			},
			email:       "test@example.com",
			password:    "password",
			expectedErr: apperr.NewApplicationError(apperr.ErrBadRequest, apperr.ErrInvalidData.Error(), apperr.ErrInvalidData),
		},
		"ユーザーを作成できない場合": {
			setupMock: func() {
				userService.EXPECT().IsEmailAvailable(gomock.Any()).Return(true, nil)
				authenticator.EXPECT().SignUp(gomock.Any(), "test@example.com", "password").Return("", apperr.ErrInvalidData)
			},
			email:       "test@example.com",
			password:    "password",
			expectedErr: apperr.NewApplicationError(apperr.ErrBadRequest, apperr.ErrInvalidData.Error(), apperr.ErrInvalidData),
		},
		"ユーザーをDBに作成できない場合": {
			setupMock: func() {
				userService.EXPECT().IsEmailAvailable(gomock.Any()).Return(true, nil)
				authenticator.EXPECT().SignUp(gomock.Any(), "test@example.com", "password").Return("user-id", nil)
				userRepo.EXPECT().Create(gomock.Any(), gomock.AssignableToTypeOf(&model.User{})).Return(errors.New("create failed"))
			},
			email:       "test@example.com",
			password:    "password",
			expectedErr: errors.New("create failed"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.setupMock()
			err := usecase.Execute(context.TODO(), tt.email, tt.password)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
