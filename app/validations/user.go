package validations

// UserUpdateRequest
// @Description User update request
type UserUpdateRequest struct {
	Name string `json:"name,omitempty" validate:"omitempty,min=2,max=50"`
}

// UserProfileUpdateRequest
// @Description User profile update request
type UserProfileUpdateRequest struct {
	Picture   string `json:"picture,omitempty" validate:"omitempty,url"`
	Gender    string `json:"gender,omitempty" validate:"oneof=F M O,omitempty,max=1"`
	Birthdate string `json:"birthdate,omitempty" validate:"omitempty"`
	Education string `json:"education,omitempty" validate:"omitempty"`
	Locale    string `json:"locale,omitempty" validate:"omitempty"`
	Timezone  string `json:"timezone,omitempty" validate:"omitempty"`
}

// ChangePasswordRequest is the request body for updating a user's password
type ChangePasswordRequest struct {
	CurrentPassword string `validate:"required,password" json:"current_password"`
	NewPassword     string `validate:"required,password" json:"new_password"`
}

// ChangeEmailRequest is the request body for updating a user's email
type ChangeEmailRequest struct {
	Email    string `validate:"required,email,max=64" json:"email"`
	Password string `validate:"required,password" json:"password"`
}

// UserPhoneNumberUpdateRequest is the request body for updating a user's phone number
type UserPhoneNumberUpdateRequest struct {
	PhoneNumber string `validate:"required" json:"phone_number"`
}

// VerifyRequest is the request body for verifying a user's email
type VerifyRequest struct {
	Code string `validate:"required" json:"code"`
}

// SendVerificationCodeRequest is the request body for sending a verification email or phone
type SendVerificationCodeRequest struct {
	Type string `validate:"required" json:"type"`
}
