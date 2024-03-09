package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/cristalhq/jwt/v4"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	v1 "github.com/usercoredev/proto/api/v1"
	"github.com/usercoredev/usercore/app/responses"
	"github.com/usercoredev/usercore/app/validations"
	"github.com/usercoredev/usercore/cache"
	"github.com/usercoredev/usercore/database"
	"github.com/usercoredev/usercore/utils"
	"github.com/usercoredev/usercore/utils/server"
	"github.com/usercoredev/usercore/utils/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"os"
	"time"
)

type UserServer struct {
	server.AuthorizationRequired
	v1.UnimplementedUserServiceServer
}

func (s *UserServer) IsAuthorizationRequired() bool {
	return true
}

func (s *UserServer) VerifyToken(ctx context.Context, in *v1.VerifyTokenRequest) (*v1.AuthenticationResponse, error) {
	claims := ctx.Value("claims").(*token.Token)
	refreshToken := in.RefreshToken

	if refreshToken == "" {
		return nil, status.Errorf(codes.PermissionDenied, responses.InvalidToken)
	}

	session, err := database.GetSessionByRefreshToken(refreshToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.PermissionDenied, responses.InvalidToken)
		}
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}

	if !session.IsActive() {
		return nil, status.Errorf(codes.PermissionDenied, responses.InvalidToken)
	}

	if !session.SessionBelongsToUser(uuid.MustParse(claims.ID)) {
		return nil, status.Errorf(codes.PermissionDenied, responses.InvalidToken)
	}

	response, err := session.RefreshUserToken()
	if err != nil {
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}

	return &v1.AuthenticationResponse{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
		ExpiresIn:    response.ExpiresIn,
	}, nil
}

func (s *UserServer) GetUser(ctx context.Context, _ *v1.GetUserRequest) (*v1.GetUserResponse, error) {
	claims := ctx.Value("claims").(*token.Token)
	userCacheKey := fmt.Sprintf(cache.UserCacheKey, claims.ID)

	var cachedUser database.User
	cacheErr := cache.Get(userCacheKey, &cachedUser)
	if cacheErr != nil {
		if !errors.Is(cacheErr, redis.Nil) && !errors.Is(cacheErr, cache.NotEnabledErr) {
			return nil, status.Errorf(codes.NotFound, responses.NotFound)
		}
	}
	if cacheErr == nil {
		return userToGetUserResponse(&cachedUser), nil
	}

	user, err := database.GetUserByID(uuid.MustParse(claims.ID), false)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, responses.NotFound)
		}
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}

	err = cache.Set(userCacheKey, user, cache.UserCacheExpiration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}
	return userToGetUserResponse(user), nil
}

func userToGetUserResponse(user *database.User) *v1.GetUserResponse {
	return &v1.GetUserResponse{
		User: &v1.User{
			Id:        user.ID.String(),
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: timestamppb.New(user.CreatedAt).AsTime().String(),
			UpdatedAt: timestamppb.New(user.UpdatedAt).AsTime().String(),
		},
	}
}

func (s *UserServer) GetUserProfile(ctx context.Context, _ *v1.GetUserProfileRequest) (*v1.GetUserProfileResponse, error) {
	claims := ctx.Value("claims").(*token.Token)
	userProfileCacheKey := fmt.Sprintf(cache.UserProfileCacheKey, claims.ID)

	var cachedProfile database.Profile
	cacheErr := cache.Get(userProfileCacheKey, &cachedProfile)
	if cacheErr != nil {
		if !errors.Is(cacheErr, redis.Nil) {
			return nil, status.Errorf(codes.NotFound, responses.ProfileNotFound)
		}
	}
	if cacheErr == nil {
		return profileToGetUserProfileResponse(&cachedProfile), nil
	}

	user, err := database.GetUserProfile(uuid.MustParse(claims.ID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, responses.NotFound)
		}
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}

	if user.Profile == nil {
		return nil, status.Errorf(codes.NotFound, responses.ProfileNotFound)
	}

	err = cache.Set(userProfileCacheKey, user.Profile, cache.UserProfileCacheExpiration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}
	return profileToGetUserProfileResponse(user.Profile), nil
}

func profileToGetUserProfileResponse(profile *database.Profile) *v1.GetUserProfileResponse {
	var birthdate string

	if profile.Birthdate != nil {
		birthdate = profile.Birthdate.Format(os.Getenv("BIRTHDATE_LAYOUT"))
	}

	return &v1.GetUserProfileResponse{
		Profile: &v1.Profile{
			Birthdate: &birthdate,
			Picture:   &profile.Picture,
			Education: &profile.Education,
			Gender:    &profile.Gender,
			Locale:    &profile.Locale,
			Timezone:  &profile.Timezone,
		},
	}
}

// UpdateUser TODO: Refactor this function
func (s *UserServer) UpdateUser(ctx context.Context, in *v1.UpdateUserRequest) (*v1.DefaultResponse, error) {
	claims := ctx.Value("claims").(*token.Token)

	userCacheKey := fmt.Sprintf(cache.UserCacheKey, claims.ID)
	userProfileCacheKey := fmt.Sprintf(cache.UserProfileCacheKey, claims.ID)

	user, err := database.GetUserByID(uuid.MustParse(claims.ID), true)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, responses.NotFound)
		}
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}

	if in.Profile != nil {
		birthdate := time.Time{}
		if in.Profile.Birthdate != nil {
			birthdate, err = time.Parse(os.Getenv("BIRTHDATE_LAYOUT"), *in.Profile.Birthdate)
			if err != nil {
				return nil, status.Errorf(codes.Internal, responses.ServerError)
			}
		}
		if user.Profile == nil {
			user.Profile = &database.Profile{
				Picture:   *in.Profile.Picture,
				Birthdate: &birthdate,
				Education: *in.Profile.Education,
				Gender:    *in.Profile.Gender,
				Locale:    *in.Profile.Locale,
				Timezone:  *in.Profile.Timezone,
			}
		} else {
			if in.Profile.Education != nil {
				user.Profile.Education = *in.Profile.Education
			}
			if in.Profile.Birthdate != nil {
				user.Profile.Birthdate = &birthdate
			}
			if in.Profile.Picture != nil {
				user.Profile.Picture = *in.Profile.Picture
			}
			if in.Profile.Gender != nil {
				user.Profile.Gender = *in.Profile.Gender
			}
			if in.Profile.Locale != nil {
				user.Profile.Locale = *in.Profile.Locale
			}
			if in.Profile.Timezone != nil {
				user.Profile.Timezone = *in.Profile.Timezone
			}
		}
	}

	if in.Name != nil {
		user.Name = *in.Name
	}

	if err := database.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&user).Error; err != nil {
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}

	if err := cache.Set(userCacheKey, user, cache.UserCacheExpiration); err != nil {
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}
	if user.Profile != nil {
		if err := cache.Set(userProfileCacheKey, user.Profile, cache.UserProfileCacheExpiration); err != nil {
			return nil, status.Errorf(codes.Internal, responses.ServerError)
		}
	}

	return &v1.DefaultResponse{
		Success: true,
	}, nil
}

func (s *UserServer) DeleteUser(ctx context.Context, _ *v1.DeleteUserRequest) (*v1.DefaultResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, responses.NotImplemented)
}

func (s *UserServer) ChangeEmail(ctx context.Context, in *v1.ChangeEmailRequest) (*v1.DefaultResponse, error) {
	claims := ctx.Value("claims").(*token.Token)

	userCacheKey := fmt.Sprintf(cache.UserCacheKey, claims.ID)

	changeEmailRequest := &validations.ChangeEmailRequest{
		Email:    in.Email,
		Password: in.Password,
	}
	validationErr := validations.ValidateStruct(changeEmailRequest)
	if validationErr != nil {
		return nil, status.Errorf(codes.InvalidArgument, responses.ValidationError)
	}

	user, err := database.GetUserByID(uuid.MustParse(claims.ID), false)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, responses.NotFound)
		}
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}

	if user.Email == changeEmailRequest.Email {
		return nil, status.Errorf(codes.InvalidArgument, responses.UserExists)
	}

	if !user.ComparePassword(changeEmailRequest.Password) {
		return nil, status.Errorf(codes.InvalidArgument, responses.InvalidCredentials)
	}

	user.UpdateUserEmail(changeEmailRequest.Email)
	if err = database.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&user).Error; err != nil {
		if utils.SQLDuplicateError(err) {
			return nil, status.Errorf(codes.AlreadyExists, responses.UserExists)
		}
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}

	if err := cache.Set(userCacheKey, user, cache.UserCacheExpiration); err != nil {
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}

	return &v1.DefaultResponse{
		Success: true,
	}, nil
}

func (s *UserServer) ChangePassword(ctx context.Context, in *v1.ChangePasswordRequest) (*v1.DefaultResponse, error) {
	claims := ctx.Value("claims").(*jwt.RegisteredClaims)

	changePasswordRequest := &validations.ChangePasswordRequest{
		CurrentPassword: in.OldPassword,
		NewPassword:     in.NewPassword,
	}

	validationErr := validations.ValidateStruct(changePasswordRequest)
	if validationErr != nil {
		return nil, status.Errorf(codes.InvalidArgument, responses.ValidationError)
	}

	user, err := database.GetUserByID(uuid.MustParse(claims.ID), false)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, responses.NotFound)
		}
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}

	if !user.ComparePassword(changePasswordRequest.CurrentPassword) {
		return nil, status.Errorf(codes.InvalidArgument, responses.InvalidCredentials)
	}

	err = user.SetPassword(changePasswordRequest.NewPassword)
	if err != nil {
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}

	if err = database.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&user).Error; err != nil {
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}

	return &v1.DefaultResponse{
		Success: true,
	}, nil
}

func (s *UserServer) SendVerificationCode(ctx context.Context, in *v1.SendVerificationCodeRequest) (*v1.DefaultResponse, error) {
	claims := ctx.Value("claims").(*jwt.RegisteredClaims)

	user, err := database.GetUserByID(uuid.MustParse(claims.ID), false)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, responses.NotFound)
		}
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}

	if in.Type == v1.VerificationType_EMAIL {
		if user.EmailVerified {
			return nil, status.Errorf(codes.InvalidArgument, responses.AlreadyVerified)
		}

		if !utils.CompareTimesByGivenMinute(utils.GetCurrentTime(), user.EmailVerifySentAt, 3) {
			return nil, status.Errorf(codes.ResourceExhausted, responses.TooManyVerifyRequest)
		}

		otpCode, err := utils.GenerateOTPCode()
		// TODO: Remove this line
		fmt.Println(otpCode)

		if err != nil {
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
			return nil, status.Errorf(codes.Aborted, responses.ServerError)
		}

		/*
				go email.Send(isEmailSent, otpCode, user.Email, language, user.Name, "Verify your email address")

				if <-isEmailSent {
					user.SetEmailVerifyCode(otpCode)
					if err := database.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&user).Error; err != nil {
						return nil, status.Errorf(codes.Internal, responses.ServerError)
					}
					return &v1.DefaultResponse{
						Success: true,
					}, nil
				}
				return nil, status.Errorf(codes.Internal, responses.ServerError)


			return nil, status.Errorf(codes.Aborted, responses.NotSupported)
		*/
	}
	return nil, status.Errorf(codes.Unimplemented, responses.NotImplemented)
}

func (s *UserServer) Verify(ctx context.Context, in *v1.VerifyRequest) (*v1.DefaultResponse, error) {
	claims := ctx.Value("claims").(*jwt.RegisteredClaims)
	userCacheKey := fmt.Sprintf(cache.UserCacheKey, claims.ID)

	verifyRequest := &validations.VerifyRequest{
		Code: in.Code,
	}

	validationErr := validations.ValidateStruct(verifyRequest)
	if validationErr != nil {
		return nil, status.Errorf(codes.InvalidArgument, responses.ValidationError)
	}

	user, err := database.GetUserByID(uuid.MustParse(claims.ID), false)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, responses.NotFound)
		}
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}

	if user.EmailVerified {
		return nil, status.Errorf(codes.InvalidArgument, responses.AlreadyVerified)
	}

	if in.Type == v1.VerificationType_EMAIL {
		if user.VerifyEmail(in.Code) {
			if err = database.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&user).Error; err != nil {
				return nil, status.Errorf(codes.Internal, responses.ServerError)
			}
			if err := cache.Set(userCacheKey, user, cache.UserCacheExpiration); err != nil {
				return nil, status.Errorf(codes.Internal, responses.ServerError)
			}
			return &v1.DefaultResponse{
				Success: true,
			}, nil
		}
		return nil, status.Errorf(codes.InvalidArgument, responses.InvalidCode)
	}

	return nil, status.Errorf(codes.Aborted, responses.NotSupported)
}

func (s *UserServer) GetUsers(ctx context.Context, in *v1.ListRequest) (*v1.GetUsersResponse, error) {
	// TODO: Implement better role system
	if !ctx.Value("claims").(*token.Token).HasRole("admin") {
		return nil, status.Errorf(codes.PermissionDenied, responses.Forbidden)
	}

	md := utils.PageMetadata{
		OrderBy:  in.OrderBy,
		Order:    in.Order,
		PageSize: in.PageSize,
		Search:   in.Search,
		Page:     in.Page,
	}

	userToSessionsResponse := func(sessions []database.Session) []*v1.Session {
		var sessionsResponse []*v1.Session
		for _, session := range sessions {
			sessionsResponse = append(sessionsResponse, &v1.Session{
				Id:         session.ID,
				ClientId:   session.ClientID,
				ClientName: session.ClientName,
				ExpiresAt:  timestamppb.New(session.ExpiresAt).AsTime().String(),
			})
		}
		return sessionsResponse
	}

	userToProfileResponse := func(profile *database.Profile) *v1.Profile {
		if profile == nil {
			return nil
		}
		var birthdate string

		if profile.Birthdate != nil {
			birthdate = profile.Birthdate.Format(os.Getenv("BIRTHDATE_LAYOUT"))
		}
		return &v1.Profile{
			Birthdate: &birthdate,
			Picture:   &profile.Picture,
			Education: &profile.Education,
			Gender:    &profile.Gender,
			Locale:    &profile.Locale,
			Timezone:  &profile.Timezone,
		}
	}

	transformResponse := func(users []*database.User) []*v1.User {
		var usersResponse []*v1.User
		for _, user := range users {
			userProfile := userToProfileResponse(user.Profile)
			userSessions := userToSessionsResponse(user.Sessions)
			usersResponse = append(usersResponse, &v1.User{
				Id:        user.ID.String(),
				Name:      user.Name,
				Email:     user.Email,
				CreatedAt: timestamppb.New(user.CreatedAt).AsTime().String(),
				UpdatedAt: timestamppb.New(user.UpdatedAt).AsTime().String(),
				Sessions:  userSessions,
				Profile:   userProfile,
			})
		}
		return usersResponse
	}

	users, count, err := database.GetUsers(md)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, responses.NotFound)
		}
		return nil, status.Errorf(codes.Internal, responses.ServerError)
	}
	var usersResponse = transformResponse(users)

	md.SetTotalCount(int32(count))
	md.SetPage(in.Page)
	return &v1.GetUsersResponse{
		Users: usersResponse,
		Meta: &v1.Meta{
			Page:       md.Page,
			TotalCount: md.TotalCount,
			TotalPages: md.TotalPages,
			PageSize:   md.PageSize,
			HasNext:    md.HasNext,
			HasPrev:    md.HasPrev,
			OrderBy:    md.OrderBy,
			Order:      md.Order,
		},
	}, nil
}
