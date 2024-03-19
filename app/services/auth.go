package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	v1 "github.com/usercoredev/proto/api/v1"
	"github.com/usercoredev/usercore/app/responses"
	"github.com/usercoredev/usercore/app/validations"
	"github.com/usercoredev/usercore/internal/client"
	database2 "github.com/usercoredev/usercore/internal/database"
	"github.com/usercoredev/usercore/internal/dateutil"
	"github.com/usercoredev/usercore/internal/textutil"
	"github.com/usercoredev/usercore/internal/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"time"
)

type AuthenticationServer struct {
	token.AuthorizationRequired
	v1.UnimplementedAuthenticationServiceServer
}

func (s *AuthenticationServer) IsAuthorizationRequired() bool {
	return false
}

func (s *AuthenticationServer) SignUp(ctx context.Context, in *v1.SignUpRequest) (*v1.AuthenticationResponse, error) {
	signUpRequest := validations.SignUpRequest{
		Email:    in.Email,
		Password: in.Password,
		Name:     in.Name,
	}
	validationErr := validations.ValidateStruct(signUpRequest)
	if validationErr != nil {
		return nil, status.Errorf(codes.InvalidArgument, responses.ValidationError)
	}

	user, err := database2.GetUserByEmail(signUpRequest.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}

	if user != nil {
		return nil, status.Errorf(codes.AlreadyExists, responses.UserExists)
	}

	var newUser = database2.User{
		Name:  signUpRequest.Name,
		Email: signUpRequest.Email,
	}
	err = newUser.SetPassword(signUpRequest.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}
	newUser.ID = uuid.New()

	if err := database2.DB.Model(&database2.User{}).Create(&newUser).Error; err != nil {
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}

	result, err := newUser.CreateSession(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}

	return &v1.AuthenticationResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}, nil
}

func (s *AuthenticationServer) SignIn(ctx context.Context, in *v1.SignInRequest) (*v1.AuthenticationResponse, error) {
	signInRequest := validations.SignInRequest{
		Email:    in.Email,
		Password: in.Password,
	}
	validationErr := validations.ValidateStruct(signInRequest)
	if validationErr != nil {
		return nil, status.Errorf(codes.InvalidArgument, responses.ValidationError)
	}

	user, err := database2.GetUserByEmail(signInRequest.Email)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, responses.InvalidCredentials)
	}

	if !user.ComparePassword(signInRequest.Password) {
		return nil, status.Errorf(codes.InvalidArgument, responses.InvalidCredentials)
	}

	result, err := user.CreateSession(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}

	return &v1.AuthenticationResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}, nil
}

func (s *AuthenticationServer) RefreshToken(ctx context.Context, in *v1.RefreshTokenRequest) (*v1.AuthenticationResponse, error) {
	ctxClient := ctx.Value(client.Key).(*client.Item)

	refreshTokenRequest := validations.RefreshTokenRequest{
		RefreshToken: in.RefreshToken,
	}

	validationErr := validations.ValidateStruct(refreshTokenRequest)
	if validationErr != nil {
		return nil, status.Errorf(codes.InvalidArgument, responses.ValidationError)
	}

	session, err := database2.GetSessionByRefreshToken(refreshTokenRequest.RefreshToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, responses.SessionNotFound)
		}
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}

	if session.ExpiresAt.Before(time.Now()) {
		return nil, status.Errorf(codes.PermissionDenied, responses.SessionExpired)
	}

	if session.ClientID != ctxClient.ID {
		return nil, status.Errorf(codes.PermissionDenied, responses.InvalidClient)
	}

	response, err := session.RefreshUserToken()
	if err != nil {
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}

	return &v1.AuthenticationResponse{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
	}, nil

}

func (s *AuthenticationServer) ResetPassword(_ context.Context, in *v1.ResetPasswordRequest) (*v1.ResetPasswordResponse, error) {
	resetPasswordRequest := validations.ResetPasswordRequest{
		Email: in.Email,
	}
	validationErr := validations.ValidateStruct(resetPasswordRequest)
	if validationErr != nil {
		return nil, status.Errorf(codes.InvalidArgument, responses.ValidationError)
	}

	otpCode := textutil.GenerateOTPCode()
	if otpCode == "" {
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}

	// TODO: Remove this line
	fmt.Println(otpCode)

	user, err := database2.GetUserByEmail(resetPasswordRequest.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.Aborted, responses.InvalidCredentials)
		}
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}

	lastReset, err := user.GetLastPasswordReset()
	if err != nil {
		lastReset = &database2.PasswordReset{
			UserID: user.ID,
		}
	}
	if !dateutil.CompareTimesByGivenMinute(dateutil.GetCurrentTime(), &lastReset.CreatedAt, 15) {
		return nil, status.Errorf(codes.ResourceExhausted, responses.TooManyResetRequest)
	}

	// TODO: Implement email sending
	/*
		language := os.Getenv("APP_DEFAULT_LANGUAGE")
		if user.Profile != nil && len(user.Profile.Locale) > 0 {
			language = utils.GetLanguage(user.Profile.Locale)
		}
		isEmailSent := make(chan bool)
	*/
	if user.Email == "" {
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}

	return nil, status.Errorf(codes.Unimplemented, responses.NotImplemented)
	/*
		go email.Send(isEmailSent, otpCode, user.Email, language, user.Name, notification.PasswordReset)

		if <-isEmailSent {
			lastReset.ResetToken = otpCode

			if err = database.DB.Model(&database.PasswordReset{}).Create(&lastReset).Error; err != nil {
				return nil, status.Errorf(codes.Internal, responses.ServerError)
			}

			return &v1.ResetPasswordResponse{
				Email: user.Email,
			}, nil
		} else {
			return nil, status.Errorf(codes.Internal, responses.ServerError)
		}
	*/
}

func (s *AuthenticationServer) ResetPasswordConfirm(_ context.Context, in *v1.ResetPasswordConfirmRequest) (*v1.DefaultResponse, error) {
	resetPasswordConfirmRequest := validations.ResetPasswordCompleteRequest{
		Email:    in.Email,
		Password: in.Password,
		Token:    in.Token,
	}
	validationErr := validations.ValidateStruct(resetPasswordConfirmRequest)
	if validationErr != nil {
		return nil, status.Errorf(codes.InvalidArgument, responses.ValidationError)
	}

	user, err := database2.GetUserByEmail(resetPasswordConfirmRequest.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.Aborted, responses.InvalidCredentials)
		}
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}

	lastReset, err := user.GetLastPasswordReset()
	if err != nil {
		return nil, status.Errorf(codes.Aborted, responses.InvalidCode)
	}

	if dateutil.CompareTimesByGivenDay(dateutil.GetCurrentTime(), &lastReset.CreatedAt, 1) {
		return nil, status.Errorf(codes.Aborted, responses.CodeExpired)
	}

	if !user.CheckPasswordResetCode(resetPasswordConfirmRequest.Token) {
		return nil, status.Errorf(codes.Aborted, responses.InvalidCode)
	}

	if err := user.SetPassword(resetPasswordConfirmRequest.Password); err != nil {
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}
	if err = database2.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&user).Error; err != nil {
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}
	/*

		language := os.Getenv("APP_DEFAULT_LANGUAGE")
		if user.Profile != nil && len(user.Profile.Locale) > 0 {
			language = utils.GetLanguage(user.Profile.Locale)
		}

		isEmailSent := make(chan bool)
	*/

	if user.Email == "" {
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}

	return nil, status.Errorf(codes.Unimplemented, responses.NotImplemented)

	/*
		go email.Send(isEmailSent, "", user.Email, language, user.Name, notification.PasswordChanged)

		if <-isEmailSent {

			if err = database.DB.Delete(&lastReset).Error; err != nil {
				return nil, status.Errorf(codes.Internal, responses.ServerError)
			}

			return &v1.DefaultResponse{
				Success: true,
			}, nil
		} else {
			return nil, status.Errorf(codes.Internal, responses.ServerError)
		}
	*/

}
