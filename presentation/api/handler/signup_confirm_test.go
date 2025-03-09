package handler_test

import (
	"net/http"
	"net/http/httptest"
	"share-basket-auth-service/core/apperr"
	"share-basket-auth-service/presentation/api/handler"
	. "share-basket-auth-service/tests/mock/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSignUpConfirmHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	usecase := NewMockSignUpConfirmUseCase(ctrl)

	defer ctrl.Finish()

	tests := map[string]struct {
		setupMock      func()
		reqBody        map[string]interface{}
		expectedStatus int
	}{
		"正常系": {
			setupMock: func() {
				usecase.EXPECT().Execute(gomock.Any(), "test@example.com", "1234567").Return(nil)
			},
			reqBody: map[string]interface{}{
				"email":             "test@example.com",
				"confirmation_code": "1234567",
			},
			expectedStatus: http.StatusOK,
		},
		"会員登録に失敗した場合": {
			setupMock: func() {
				usecase.EXPECT().Execute(gomock.Any(), "test@example.com", "1234567").Return(apperr.NewApplicationError(apperr.ErrBadRequest, "invalid input", apperr.ErrInvalidData))
			},
			reqBody: map[string]interface{}{
				"email":             "test@example.com",
				"confirmation_code": "1234567",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.setupMock()

			req := newJSONRequest(http.MethodPost, "/signup/confirm", tt.reqBody)
			rec := httptest.NewRecorder()
			s := setupTestServer()

			s.POST("/signup/confirm", handler.MakeSignUpConfirmHandler(usecase))
			s.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
