package errorutil

import (
	"fmt"
)

type UCError struct {
	Code    int
	Message string
}

var (
	ErrGRPCPortNotSet         = &UCError{Code: 1000, Message: "gRPC Server port not set"}
	ErrGRPCFailedToServe      = &UCError{Code: 1000, Message: "Failed to serve gRPC server"}
	ErrClientFilePathNotSet   = &UCError{Code: 1001, Message: "Client file path not set"}
	ErrPrivateKeyPathNotSet   = &UCError{Code: 1002, Message: "Private key path not set"}
	ErrPublicKeyPathNotSet    = &UCError{Code: 1003, Message: "Public key path not set"}
	ErrNoClients              = &UCError{Code: 1001, Message: "No clients found"}
	ErrInvalidRequest         = &UCError{Code: 1002, Message: "Invalid request"}
	ErrValidationError        = &UCError{Code: 1003, Message: "Validation error"}
	ErrServerError            = &UCError{Code: 1004, Message: "Server error"}
	ErrInvalidCredentials     = &UCError{Code: 1005, Message: "Invalid credentials"}
	ErrClientRequired         = &UCError{Code: 1006, Message: "Client required"}
	ErrInvalidClient          = &UCError{Code: 1007, Message: "Invalid client"}
	ErrNotImplemented         = &UCError{Code: 1008, Message: "Not implemented"}
	ErrForbidden              = &UCError{Code: 1009, Message: "Forbidden"}
	ErrAlreadyExists          = &UCError{Code: 1010, Message: "Already exists"}
	ErrNotFound               = &UCError{Code: 1011, Message: "Not found"}
	ErrProfileNotFound        = &UCError{Code: 1012, Message: "Profile not found"}
	ErrSessionNotFound        = &UCError{Code: 1013, Message: "Session not found"}
	ErrSessionExpired         = &UCError{Code: 1014, Message: "Session expired"}
	ErrTooManyVerifyRequest   = &UCError{Code: 1015, Message: "Too many verify request"}
	ErrTooManyResetRequest    = &UCError{Code: 1016, Message: "Too many reset request"}
	ErrInvalidCode            = &UCError{Code: 1017, Message: "Invalid reset code"}
	ErrCodeExpired            = &UCError{Code: 1018, Message: "Reset code expired"}
	ErrAlreadyVerified        = &UCError{Code: 1019, Message: "Already verified"}
	ErrSocialProviderNotFound = &UCError{Code: 1021, Message: "Social provider not found"}
	ErrResetPasswordCodeSent  = &UCError{Code: 1022, Message: "Reset password code sent"}
	ErrTokenExpired           = &UCError{Code: 1023, Message: "Token expired"}
	ErrInvalidToken           = &UCError{Code: 1024, Message: "Invalid token"}
	ErrTokenMalformed         = &UCError{Code: 1025, Message: "Token malformed"}
	ErrNotSupported           = &UCError{Code: 1026, Message: "Not supported"}
)

func (e *UCError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}
