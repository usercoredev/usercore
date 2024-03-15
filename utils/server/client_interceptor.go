package server

import (
	"context"
	"github.com/usercoredev/usercore/app/responses"
	"github.com/usercoredev/usercore/utils/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func ClientInterceptor(clients []client.Item) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok || len(md.Get("client")) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, responses.ClientRequired)
		}
		clientID := md.Get("client")[0]
		if clientID == "" {
			return nil, status.Errorf(codes.Unauthenticated, responses.ClientRequired)
		}
		mdClient := client.GetClient(clientID, clients)
		if mdClient == nil {
			return nil, status.Errorf(codes.Unauthenticated, responses.InvalidClient)
		}
		ctx = context.WithValue(ctx, client.Key, mdClient)
		return handler(ctx, req)
	}
}
