package services

import (
	v1 "github.com/usercoredev/proto/api/v1"
	"github.com/usercoredev/usercore/utils/server"
)

type RoleServer struct {
	server.AuthorizationRequired
	v1.UnimplementedRoleServiceServer
}

func (s *RoleServer) IsAuthorizationRequired() bool {
	return true
}
