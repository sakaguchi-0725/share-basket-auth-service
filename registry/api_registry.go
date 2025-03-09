package registry

import (
	"context"
	"fmt"
	"share-basket-auth-service/core/config"
	"share-basket-auth-service/domain/service"
	"share-basket-auth-service/infrastructure/persistence"
	"share-basket-auth-service/presentation/api/handler"
	"share-basket-auth-service/usecase"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type apiRegistry struct {
	signUpUseCase        usecase.SignUpUseCase
	signUpConfirmUseCase usecase.SignUpConfirmUseCase
}

func InjectAPI(ctx context.Context, db *gorm.DB, awsConfig config.AWSConfig) (*apiRegistry, error) {
	authenticator, err := persistence.NewCognito(ctx, awsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create authenticator: %w", err)
	}
	userRepo := persistence.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	return &apiRegistry{
		signUpUseCase:        usecase.NewSignUpUseCase(authenticator, userRepo, userService),
		signUpConfirmUseCase: usecase.NewSignUpConfirmUseCase(authenticator),
	}, nil
}

func (r *apiRegistry) SignUpHandler() echo.HandlerFunc {
	return handler.MakeSignUpHandler(r.signUpUseCase)
}

func (r *apiRegistry) SignUpConfirmHandler() echo.HandlerFunc {
	return handler.MakeSignUpConfirmHandler(r.signUpConfirmUseCase)
}
