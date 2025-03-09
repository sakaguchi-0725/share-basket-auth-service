package registry

import (
	"context"
	"fmt"
	"share-basket-auth-service/core/config"
	"share-basket-auth-service/gen"
	"share-basket-auth-service/infrastructure/persistence"
	"share-basket-auth-service/presentation/grpc/handler"
	"share-basket-auth-service/usecase"

	"gorm.io/gorm"
)

type gRPCRegistry struct {
	verifyTokenUseCase usecase.VerifyTokenUseCase
}

func InjectGRPC(ctx context.Context, db *gorm.DB, awsConfig config.AWSConfig) (*gRPCRegistry, error) {
	authenticator, err := persistence.NewCognito(ctx, awsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create authenticator: %w", err)
	}
	userRepo := persistence.NewUserRepository(db)

	return &gRPCRegistry{
		verifyTokenUseCase: usecase.NewVerifyTokenUseCase(authenticator, userRepo),
	}, nil
}

func (r *gRPCRegistry) VerifyTokenHandler() gen.AuthServiceServer {
	return handler.NewVerifyTokenHandler(r.verifyTokenUseCase)
}
