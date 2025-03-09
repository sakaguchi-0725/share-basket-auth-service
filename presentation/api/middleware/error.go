package middleware

import (
	"net/http"
	"share-basket-auth-service/core/apperr"

	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func ErrorMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}

			if appErr, ok := err.(apperr.ApplicationError); ok {
				errCode := appErr.Code()

				switch errCode {
				case apperr.ErrBadRequest:
					return c.JSON(http.StatusBadRequest, &ErrorResponse{
						Code:    errCode.String(),
						Message: appErr.Message(),
					})
				case apperr.ErrNotFound:
					return c.JSON(http.StatusNotFound, &ErrorResponse{
						Code:    errCode.String(),
						Message: appErr.Message(),
					})
				case apperr.ErrUnauthorized:
					return c.JSON(http.StatusUnauthorized, &ErrorResponse{
						Code:    errCode.String(),
						Message: appErr.Message(),
					})
				}
			}

			return c.JSON(http.StatusInternalServerError, &ErrorResponse{
				Code:    "InternalServerError",
				Message: "サーバーでエラーが発生しました",
			})
		}
	}
}
