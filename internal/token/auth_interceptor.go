package token

import (
	"context"
	"errors"
	"github.com/cristalhq/jwt/v4"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/usercoredev/usercore/app/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"time"
)

type AuthorizationRequired interface {
	IsAuthorizationRequired() bool
}

func (s *Settings) AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if authConfigProvider, ok := info.Server.(AuthorizationRequired); ok {
			if authConfigProvider.IsAuthorizationRequired() {
				_, ok = metadata.FromIncomingContext(ctx)
				if !ok {
					return nil, status.Errorf(codes.Unauthenticated, "metadata not provided")
				}
				bearerToken, err := auth.AuthFromMD(ctx, s.Scheme)
				if err != nil {
					return nil, err
				}
				claims, err := s.verify(bearerToken)
				if err != nil {
					return nil, status.Errorf(codes.Unauthenticated, err.Error())
				}
				ctx = context.WithValue(ctx, Claims, claims)
			}
		}
		return handler(ctx, req)
	}
}

func (s *Settings) verify(receivedToken string) (jwt.RegisteredClaims, error) {
	var newClaims jwt.RegisteredClaims
	err := jwt.ParseClaims([]byte(receivedToken), s.Verifier, &newClaims)
	if err != nil {
		return jwt.RegisteredClaims{}, errors.New(responses.TokenMalformed)
	}
	var isValid = newClaims.IsValidAt(time.Now())
	if !isValid {
		return jwt.RegisteredClaims{}, errors.New(responses.TokenExpired)
	}
	return newClaims, nil
}
