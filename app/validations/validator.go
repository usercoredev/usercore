package validations

import (
	"github.com/go-playground/validator/v10"
	"github.com/talut/dotenv"
	"regexp"
	"strings"
)

type ErrorResponse struct {
	Field    string `json:"field"`
	Reason   string `json:"reason"`
	Expected string `json:"expected,omitempty"`
}

var validate = validator.New()

func ValidateStruct(validatorStruct interface{}) []ErrorResponse {
	var errors []ErrorResponse
	_ = validate.RegisterValidation("password", passwordValidator)
	err := validate.Struct(validatorStruct)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = strings.ToLower(err.Field())
			element.Reason = err.Tag()
			element.Expected = err.Param()
			errors = append(errors, element)
		}
	}
	return errors
}

func passwordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	var (
		uppercase = regexp.MustCompile(`[A-Z]`)
		lowercase = regexp.MustCompile(`[a-z]`)
		number    = regexp.MustCompile(`\d`)
		special   = regexp.MustCompile(`[^a-zA-Z\d]`)
	)
	passwordMinLength := dotenv.GetInt("PASSWORD_MIN_LENGTH", 8)
	passwordMaxLength := dotenv.GetInt("PASSWORD_MAX_LENGTH", 64)
	return len(password) >= passwordMinLength &&
		len(password) <= passwordMaxLength &&
		uppercase.MatchString(password) &&
		lowercase.MatchString(password) &&
		number.MatchString(password) &&
		special.MatchString(password)
}
