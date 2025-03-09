package handler

import (
	"net/http"
	"share-basket-auth-service/core/apperr"
	"share-basket-auth-service/usecase"

	"github.com/labstack/echo/v4"
)

type signUpConfirmRequest struct {
	Email            string `json:"email"`
	ConfirmationCode string `json:"confirmation_code"`
}

func MakeSignUpConfirmHandler(usecase usecase.SignUpConfirmUseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		req, err := newSignUpConfirmRequest(c)
		if err != nil {
			return err
		}

		err = usecase.Execute(c.Request().Context(), req.Email, req.ConfirmationCode)
		if err != nil {
			return err
		}

		return c.NoContent(http.StatusOK)
	}
}

func newSignUpConfirmRequest(c echo.Context) (signUpConfirmRequest, error) {
	var req signUpConfirmRequest

	if err := c.Bind(&req); err != nil {
		return signUpConfirmRequest{}, apperr.NewInvalidInputError(err)
	}

	return req, nil
}
