package services

import (
	v1 "github.com/usercoredev/proto/api/v1"
	"github.com/usercoredev/usercore/internal/token"
)

type RoleServer struct {
	token.AuthorizationRequired
	v1.UnimplementedRoleServiceServer
}

func (s *RoleServer) IsAuthorizationRequired() bool {
	return true
}
