package database

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/usercoredev/usercore/utils"
	"github.com/usercoredev/usercore/utils/client"
	"github.com/usercoredev/usercore/utils/token"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
	"strconv"
	"time"
)

type User struct {
	UUIDBaseModel

	Name string `gorm:"not null" json:"name"`

	Email             string     `gorm:"unique" json:"email"`
	EmailVerified     bool       `gorm:"default:false" json:"email_verified,omitempty"`
	EmailVerifyCode   string     `gorm:"default:null" json:"-"`
	EmailVerifySentAt *time.Time `gorm:"default:null" json:"-"`

	Password     string `json:"-"`
	PasswordSalt string `json:"-"`

	Sessions        []Session        `json:"sessions,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Profile         *Profile         `json:"profile,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PasswordReset   []PasswordReset  `json:"-" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SocialProviders []SocialProvider `json:"-" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.New()
	return nil
}

// ComparePassword compares the password of a user
func (user *User) ComparePassword(password string) bool {
	return utils.CheckPasswordHash(password, user.Password)
}

// SetPassword sets the password of a user
func (user *User) SetPassword(password string) error {
	hashedPassword, err := utils.GeneratePasswordHash(password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return nil
}

func (user *User) UserSessionLimiter() error {
	limit, err := strconv.Atoi(os.Getenv("MAX_SESSIONS_PER_USER"))

	if len(user.Sessions) >= limit {
		if err = DB.Delete(&user.Sessions[0]).Error; err != nil {
			return err
		}
	}
	return nil
}

// CheckPasswordResetToken checks the password reset code of a user
func (user *User) CheckPasswordResetToken(token string) bool {
	if len(user.PasswordReset) > 0 {
		return user.PasswordReset[0].CheckResetToken(token)
	}
	return false
}

// SetEmailVerifyCode verifies the email of a user
func (user *User) SetEmailVerifyCode(code string) {
	user.EmailVerifyCode = code
	currentTime := utils.GetCurrentTime()
	user.EmailVerifySentAt = &currentTime
}

// VerifyEmail verifies the email of a user
func (user *User) VerifyEmail(code string) bool {
	if user.EmailVerifyCode == code {
		user.EmailVerified = true
		user.EmailVerifyCode = ""
		user.EmailVerifySentAt = nil
		return true
	}
	return false
}

// UpdateUserEmail updates the email of a user and sets the email verified to false and email verify code to null
func (user *User) UpdateUserEmail(email string) {
	user.EmailVerified = false
	user.EmailVerifyCode = ""
	user.Email = email
}

// GetUserByEmail gets a user by email
func GetUserByEmail(e string) (*User, error) {
	if len(e) > 0 {
		var user User
		if err := userPreload().Where("email = ?", e).First(&user).Error; err != nil {
			return nil, err
		}
		return &user, nil
	}
	return nil, errors.New("email is empty")
}

func (user *User) CheckPasswordResetCode(code string) bool {
	if len(user.PasswordReset) > 0 {
		return user.PasswordReset[0].CheckResetToken(code)
	}
	return false
}

// GetUserByID gets a user by id
func GetUserByID(id uuid.UUID, preload bool) (*User, error) {
	var user User
	if !preload {
		if err := DB.Model(&User{}).First(&user, "id = ?", id).Error; err != nil {
			return nil, err
		}
	} else {
		if err := userPreload().First(&user, "id = ?", id).Error; err != nil {
			return nil, err
		}
	}
	return &user, nil
}

func GetUserProfile(id uuid.UUID) (*User, error) {
	var user User
	if err := DB.Preload("Profile").First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (user *User) GetLastPasswordReset() (*PasswordReset, error) {
	if len(user.PasswordReset) > 0 {
		return &user.PasswordReset[0], nil
	}
	return nil, errors.New("password reset not found")
}

func userPreload() *gorm.DB {
	return DB.Preload(clause.Associations)
}

func (user *User) CreateSession(ctx context.Context) (*token.DefaultToken, error) {
	sClient := ctx.Value("client").(*client.Client)
	jwt, expiresIn, err := token.CreateJWT(user.ID)
	if err != nil {
		return nil, err
	}

	rToken, refreshTokenExpireAt := token.CreateRefreshToken(user.ID)

	if err := user.UserSessionLimiter(); err != nil {
		return nil, err
	}

	var session = Session{
		UserID:       user.ID,
		RefreshToken: rToken,
		ExpiresAt:    *refreshTokenExpireAt,
		ClientID:     sClient.ID,
		ClientName:   sClient.Name,
	}
	if err := DB.Model(&Session{}).Create(&session).Error; err != nil {
		return nil, err
	}

	return &token.DefaultToken{
		AccessToken:  jwt,
		RefreshToken: rToken,
		ExpiresIn:    expiresIn,
	}, nil
}

func GetUsers(md utils.PageMetadata) ([]*User, int64, error) {
	var count int64
	var users []*User
	if err := DB.Model(User{}).
		Preload("Sessions").
		Preload("Profile").
		Joins("Profile").
		Or("name LIKE ?", "%"+md.Search+"%").
		Or("email LIKE ?", "%"+md.Search+"%").
		Or("Profile.timezone LIKE ?", "%"+md.Search+"%").
		Or("Profile.birthdate LIKE ?", "%"+md.Search+"%").
		Or("Profile.gender LIKE ?", "%"+md.Search+"%").
		Or("Profile.locale LIKE ?", "%"+md.Search+"%").
		Or("Profile.education LIKE ?", "%"+md.Search+"%").
		Count(&count).
		Order(md.ConvertToOrder()).
		Offset(int(md.Offset())).
		Limit(int(md.PageSize)).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, count, nil
}
