package handler

import (
	"net/http"
	"share-basket-auth-service/core/apperr"
	"share-basket-auth-service/usecase"

	"github.com/labstack/echo/v4"
)

type signUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func MakeSignUpHandler(usecase usecase.SignUpUseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		input, err := newSignUpRequest(c)
		if err != nil {
			return err
		}

		if err := usecase.Execute(c.Request().Context(), input.Email, input.Password); err != nil {
			return err
		}

		return c.NoContent(http.StatusOK)
	}
}

func newSignUpRequest(c echo.Context) (signUpRequest, error) {
	var req signUpRequest
	if err := c.Bind(&req); err != nil {
		return signUpRequest{}, apperr.NewInvalidInputError(err)
	}

	return req, nil
}
