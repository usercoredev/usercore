package database

import (
	"github.com/google/uuid"
)

type Device struct {
	UINTBaseModel
	UserID    uuid.UUID `gorm:"default:null" json:"-"`
	SessionID string    `gorm:"default:null" json:"session_id,omitempty"`
	Name      string    `gorm:"default:null" json:"name,omitempty"`
	IP        string    `gorm:"default:null" json:"ip,omitempty"`
	OS        string    `gorm:"default:null" json:"os,omitempty"`
	Token     string    `gorm:"default:null" json:"token,omitempty"`
}

// GetDevicesByUserId gets all devices of a user
func GetDevicesByUserId(userId uuid.UUID) ([]Device, error) {
	var devices []Device
	if err := DB.Where("user_id = ?", userId).Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}
