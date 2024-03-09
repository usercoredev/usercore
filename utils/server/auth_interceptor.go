package server

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/usercoredev/usercore/utils/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthorizationRequired interface {
	IsAuthorizationRequired() bool
}

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if authConfigProvider, ok := info.Server.(AuthorizationRequired); ok {
		if authConfigProvider.IsAuthorizationRequired() {
			_, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return nil, status.Errorf(codes.Unauthenticated, "metadata not provided")
			}

			bearerToken, err := auth.AuthFromMD(ctx, "Bearer")
			if err != nil {
				return nil, err
			}

			claims, err := token.VerifyJWT(bearerToken)
			if err != nil {
				return nil, status.Errorf(codes.Unauthenticated, err.Error())
			}
			ctx = context.WithValue(ctx, "claims", claims)
		}
	}

	// Call the handler
	return handler(ctx, req)
}
