package validations

type SignUpRequest struct {
	Name     string `validate:"required,min=3,max=64"`
	Email    string `validate:"required,email,min=5,max=64"`
	Password string `validate:"required,min=8,max=64,password"`
}

type SignInRequest struct {
	Email    string `validate:"required,email,min=5,max=64" json:"email"`
	Password string `validate:"required,min=8,max=64,password" json:"password"`
}

type StoreDeviceInfoRequest struct {
	Name  string `json:"device_name"`
	IP    string `json:"device_ip"`
	OS    string `json:"device_os"`
	Token string `json:"device_token"`
}

type SocialAuth struct {
	Provider string `validate:"required" json:"provider"`
	Name     string `validate:"max=64,omitempty" json:"name"`
	Code     string `validate:"required,min=10" json:"code"`
	Email    string `validate:"max=64,omitempty" json:"email"`
	UserId   string `validate:"omitempty" json:"user_id"`
	Picture  string `validate:"omitempty" json:"picture"`
}
type SocialAuthComplete struct {
	Name  string `validate:"required,max=64" json:"name"`
	Email string `validate:"required,email,min=5,max=64" json:"email"`
}

type SignInWithPhoneNumber struct {
	PhoneNumber string `validate:"required,min=10" json:"phone_number"`
	Password    string `validate:"required,min=8,max=64,password" json:"password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `validate:"required" json:"refresh_token"`
}

type VerifyEmailRequest struct {
	VerificationCode string `validate:"required" json:"verification_code"`
}

type VerifyPhoneRequest struct {
	VerificationCode string `validate:"required" json:"verification_code"`
}

type DeleteSessionRequest struct {
	RefreshToken string `validate:"required" json:"refresh_token"`
}

type ResetPasswordRequest struct {
	Email string `validate:"required,email,min=5,max=64" json:"email"`
}

type ResetPasswordCompleteRequest struct {
	Email    string `validate:"required,email,min=5,max=64" json:"email"`
	Token    string `validate:"required" json:"token"`
	Password string `validate:"required,min=8,max=64,password" json:"password"`
}

type ResetPasswordWithPhoneNumberRequest struct {
	PhoneNumber string `validate:"required,min=10" json:"phone_number"`
}

type ChangePasswordWithEmailRequest struct {
	Email            string `validate:"required,email,min=5,max=64" json:"email"`
	NewPassword      string `validate:"required,min=8,max=64,password" json:"new_password"`
	VerificationCode string `validate:"required" json:"verification_code"`
}

type ChangePasswordWithPhoneNumberRequest struct {
	PhoneNumber      string `validate:"required,min=10" json:"phone_number"`
	NewPassword      string `validate:"required,min=8,max=64,password" json:"new_password"`
	VerificationCode string `validate:"required" json:"verification_code"`
}
