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

func TestLoginHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	usecase := NewMockLoginUseCase(ctrl)

	defer ctrl.Finish()

	tests := map[string]struct {
		setupMock      func()
		reqBody        map[string]interface{}
		expectedCookie *http.Cookie
		expectedStatus int
	}{
		"正常系": {
			setupMock: func() {
				usecase.EXPECT().Execute(gomock.Any(), "test@example.com", "password").Return("token", nil)
			},
			reqBody: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password",
			},
			expectedCookie: &http.Cookie{
				Name:  "access_token",
				Value: "token",
			},
			expectedStatus: http.StatusOK,
		},
		"入力値が不正な場合": {
			setupMock: func() {
				usecase.EXPECT().Execute(gomock.Any(), "test@example.com", "password").Return("", apperr.NewInvalidInputError(apperr.ErrInvalidData))
			},
			reqBody: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.setupMock()

			req := newJSONRequest(http.MethodPost, "/login", tt.reqBody)
			rec := httptest.NewRecorder()
			s := setupTestServer()

			s.POST("/login", handler.MakeLoginHandler(usecase))
			s.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if rec.Code == http.StatusOK {
				cookies := rec.Result().Cookies()
				assert.Len(t, cookies, 1)

				cookie := cookies[0]
				assert.Equal(t, tt.expectedCookie.Name, cookie.Name)
				assert.Equal(t, tt.expectedCookie.Value, cookie.Value)
			}
		})
	}
}
