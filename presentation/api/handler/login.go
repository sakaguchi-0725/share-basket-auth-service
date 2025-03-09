package handler

import (
	"net/http"
	"share-basket-auth-service/core/apperr"
	"share-basket-auth-service/usecase"

	"github.com/labstack/echo/v4"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func MakeLoginHandler(usecase usecase.LoginUseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		req, err := newLoginRequest(c)
		if err != nil {
			return err
		}

		token, err := usecase.Execute(c.Request().Context(), req.Email, req.Password)
		if err != nil {
			return err
		}

		cookie := newCookie(token)
		c.SetCookie(cookie)

		return c.NoContent(http.StatusOK)
	}
}

func newLoginRequest(c echo.Context) (loginRequest, error) {
	var req loginRequest

	if err := c.Bind(&req); err != nil {
		return loginRequest{}, apperr.NewInvalidInputError(err)
	}
	return req, nil
}

func newCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:  "access_token",
		Value: token,
		Path:  "/",
		// HttpOnly: true,
		// SameSite: http.SameSiteNoneMode,
		// Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
}
