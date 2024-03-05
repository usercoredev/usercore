package database

import (
	"github.com/google/uuid"
)

type PasswordReset struct {
	UINTBaseModel
	UserID     uuid.UUID `gorm:"default:null" json:"-"`
	ResetToken string    `gorm:"default:null" json:"-"`
}

// CheckResetToken checks if the password reset token is valid
func (passwordReset *PasswordReset) CheckResetToken(token string) bool {
	return passwordReset.ResetToken == token
}
