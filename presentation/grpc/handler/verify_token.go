package handler

import (
	"context"
	"share-basket-auth-service/gen"
	"share-basket-auth-service/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type verifyTokenHandler struct {
	gen.UnimplementedAuthServiceServer
	usecase usecase.VerifyTokenUseCase
}

func NewVerifyTokenHandler(usecase usecase.VerifyTokenUseCase) gen.AuthServiceServer {
	return &verifyTokenHandler{
		usecase: usecase,
	}
}

func (h *verifyTokenHandler) VerifyToken(ctx context.Context, req *gen.VerifyTokenRequest) (*gen.VerifyTokenResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "missing metadata")
	}

	tokenValues := md.Get("access_token")
	if len(tokenValues) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing access token")
	}

	token := tokenValues[0]
	userID, err := h.usecase.Execute(ctx, token)
	if err != nil {
		// TODO: apperrで判定する
		return nil, status.Errorf(codes.Internal, "failed to verify token: %v", err)
	}

	return &gen.VerifyTokenResponse{
		UserID: userID,
	}, nil
}
