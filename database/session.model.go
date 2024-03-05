package database

import (
	"github.com/google/uuid"
	"github.com/usercoredev/usercore/utils/token"
	"time"
)

type Session struct {
	UINTBaseModel
	UserID       uuid.UUID `gorm:"default:null" json:"-"`
	RefreshToken string    `gorm:"default:null" json:"-"`
	ExpiresAt    time.Time `gorm:"default:null" json:"expires_at,omitempty"`
	ClientID     string    `gorm:"default:null" json:"client_id,omitempty"`
	ClientName   string    `gorm:"default:null" json:"client_name,omitempty"`
	Device       Device    `gorm:"foreignKey:SessionID" json:"device,omitempty"`
}

// GetSessionByRefreshToken returns a session by refresh token
func GetSessionByRefreshToken(refreshToken string) (*Session, error) {
	var session Session
	if err := DB.Where("refresh_token = ?", refreshToken).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func GetSessionById(id uint64) (*Session, error) {
	var session Session
	if err := DB.Where("id = ?", id).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func GetSessionsByUserId(id uuid.UUID) ([]Session, error) {
	var sessions []Session

	if err := DB.Model(&Session{}).Where("user_id = ?", id).Order("id desc").Find(&sessions).Error; err != nil {
		return nil, err
	}

	return sessions, nil
}

func (session *Session) SessionBelongsToUser(userId uuid.UUID) bool {
	if session.UserID != userId {
		return false
	}
	return true
}

func (session *Session) Create() error {
	if err := DB.Model(&Session{}).Create(&session).Error; err != nil {
		return err
	}
	return nil
}

func (session *Session) Delete() error {
	if err := DB.Model(&Session{}).Delete(&session).Error; err != nil {
		return err
	}
	return nil
}

func (session *Session) Update() error {
	if err := DB.Model(&Session{}).Updates(&session).Error; err != nil {
		return err
	}
	return nil
}

func (session *Session) GetDevice() (*Device, error) {
	var device Device
	if err := DB.Model(&Device{}).Where("session_id = ?", session.ID).First(&device).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

func (session *Session) GetDeviceByDeviceId(deviceId string) (*Device, error) {
	var device Device
	if err := DB.Model(&Device{}).Where("session_id = ? AND device_id = ?", session.ID, deviceId).First(&device).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

func (session *Session) GetDevices() ([]Device, error) {
	var devices []Device
	if err := DB.Model(&Device{}).Where("session_id = ?", session.ID).Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

func (session *Session) IsActive() bool {
	if session.ExpiresAt.Before(time.Now()) {
		return false
	}
	return true
}

func (session *Session) RefreshUserToken() (*token.DefaultToken, error) {

	jwt, expiresIn, err := token.CreateJWT(session.UserID)
	if err != nil {
		return nil, err
	}

	rToken, refreshTokenExpiresAt := token.CreateRefreshToken(session.UserID)
	session.RefreshToken = rToken
	session.ExpiresAt = *refreshTokenExpiresAt

	if err = DB.Save(session).Error; err != nil {
		return nil, err
	}

	return &token.DefaultToken{
		AccessToken:  jwt,
		RefreshToken: rToken,
		ExpiresIn:    expiresIn,
	}, nil
}
