package services

import (
	v1 "github.com/usercoredev/proto/api/v1"
	"github.com/usercoredev/usercore/utils/server"
)

type PermissionServer struct {
	server.AuthorizationRequired
	v1.UnimplementedPermissionServiceServer
}

func (s *PermissionServer) IsAuthorizationRequired() bool {
	return true
}
