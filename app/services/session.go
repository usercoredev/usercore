package services

import (
	"context"
	"errors"
	"github.com/cristalhq/jwt/v4"
	"github.com/google/uuid"
	v1 "github.com/usercoredev/proto/api/v1"
	"github.com/usercoredev/usercore/app/responses"
	"github.com/usercoredev/usercore/database"
	"github.com/usercoredev/usercore/utils/server"
	"github.com/usercoredev/usercore/utils/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type SessionServer struct {
	server.AuthorizationRequired
	v1.UnimplementedSessionServiceServer
}

func (s *SessionServer) IsAuthorizationRequired() bool {
	return true
}

func (s *SessionServer) GetSessions(ctx context.Context, _ *v1.GetSessionsRequest) (*v1.GetSessionsResponse, error) {
	claims := ctx.Value("claims").(*token.Token)

	userSessions, err := database.GetSessionsByUserId(uuid.MustParse(claims.ID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, responses.NotFound)
		}
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}

	var sessions []*v1.Session
	for _, session := range userSessions {
		sessions = append(sessions, &v1.Session{
			Id:         session.ID,
			ClientName: session.ClientName,
			ClientId:   session.ClientID,
			ExpiresAt:  timestamppb.New(session.ExpiresAt).AsTime().String(),
			CreatedAt:  timestamppb.New(session.CreatedAt).AsTime().String(),
			UpdatedAt:  timestamppb.New(session.UpdatedAt).AsTime().String(),
		})
	}

	return &v1.GetSessionsResponse{Sessions: sessions}, nil
}

func (s *SessionServer) DeleteSession(ctx context.Context, in *v1.DeleteSessionRequest) (*v1.DefaultResponse, error) {
	claims := ctx.Value("claims").(*jwt.RegisteredClaims)

	session, err := database.GetSessionById(in.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, responses.SessionNotFound)
		}
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}
	if session.SessionBelongsToUser(uuid.MustParse(claims.ID)) {
		if err = database.DB.Delete(&session).Error; err != nil {
			return nil, status.Errorf(codes.Internal, responses.ServerError)
		}

		return &v1.DefaultResponse{Success: true}, nil
	}
	return nil, status.Errorf(codes.PermissionDenied, responses.Forbidden)
}

func (s *SessionServer) SignOut(ctx context.Context, in *v1.SignOutRequest) (*v1.DefaultResponse, error) {
	claims := ctx.Value("claims").(*jwt.RegisteredClaims)

	session, err := database.GetSessionByRefreshToken(in.RefreshToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, responses.SessionNotFound)
		}
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}

	if session.SessionBelongsToUser(uuid.MustParse(claims.ID)) {
		if err = database.DB.Delete(&session).Error; err != nil {
			return nil, status.Errorf(codes.Internal, responses.ServerError)
		}

		return &v1.DefaultResponse{Success: true}, nil
	}

	return nil, status.Errorf(codes.PermissionDenied, responses.Forbidden)
}
