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

func TestSignUpConfirmUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	authenticator := NewMockAuthenticator(ctrl)
	usecase := usecase.NewSignUpConfirmUseCase(authenticator)

	defer ctrl.Finish()

	tests := map[string]struct {
		setupMock        func()
		email            string
		confirmationCode string
		expectedErr      error
	}{
		"正常系": {
			setupMock: func() {
				authenticator.EXPECT().SignUpConfirm(gomock.Any(), "test@example.com", "1234567").Return(nil)
			},
			email:            "test@example.com",
			confirmationCode: "1234567",
			expectedErr:      nil,
		},
		"入力値が不正な場合": {
			setupMock: func() {
				authenticator.EXPECT().SignUpConfirm(gomock.Any(), "test@example.com", "1234567").Return(apperr.ErrInvalidData)
			},
			email:            "test@example.com",
			confirmationCode: "1234567",
			expectedErr:      apperr.ErrInvalidData,
		},
		"認証コードの期限が切れている場合": {
			setupMock: func() {
				authenticator.EXPECT().SignUpConfirm(gomock.Any(), "test@example.com", "1234567").Return(apperr.ErrExpiredCodeException)
			},
			email:            "test@example.com",
			confirmationCode: "1234567",
			expectedErr:      apperr.ErrExpiredCodeException,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.setupMock()
			err := usecase.Execute(context.Background(), tt.email, tt.confirmationCode)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
