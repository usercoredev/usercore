package services

import (
	v1 "github.com/usercoredev/proto/api/v1"
	"github.com/usercoredev/usercore/utils/token"
)

type PermissionServer struct {
	token.AuthorizationRequired
	v1.UnimplementedPermissionServiceServer
}

func (s *PermissionServer) IsAuthorizationRequired() bool {
	return true
}
