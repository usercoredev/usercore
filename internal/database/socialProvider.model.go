package database

import "github.com/google/uuid"

type SocialProvider struct {
	UINTBaseModel
	UserID         uuid.UUID `gorm:"default:null" json:"-"`
	Provider       string    `gorm:"default:null" json:"provider,omitempty"`
	ProviderUserID string    `gorm:"default:null" json:"provider_user_id,omitempty"`
	AccessToken    string    `gorm:"default:null" json:"access_token,omitempty"`
	RefreshToken   string    `gorm:"default:null" json:"refresh_token,omitempty"`
}
