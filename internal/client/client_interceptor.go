package client

import (
	"context"
	"github.com/usercoredev/usercore/app/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func ClientInterceptor(clientSettings Settings) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok || len(md.Get(string(Key))) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, responses.ClientRequired)
		}
		clientID := md.Get(string(Key))[0]
		if clientID == "" {
			return nil, status.Errorf(codes.Unauthenticated, responses.ClientRequired)
		}
		mdClient := clientSettings.GetClient(clientID)
		if mdClient == nil {
			return nil, status.Errorf(codes.Unauthenticated, responses.InvalidClient)
		}
		ctx = context.WithValue(ctx, Key, mdClient)
		return handler(ctx, req)
	}
}
