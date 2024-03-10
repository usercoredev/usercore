package database

import (
	"context"
	"errors"
	"fmt"
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

	Name string `gorm:"not null" json:"name" sortable:"true"`

	Email             string     `gorm:"unique" json:"email" sortable:"true"`
	EmailVerified     bool       `gorm:"default:false" json:"email_verified,omitempty" sortable:"true"`
	EmailVerifyCode   string     `gorm:"default:null" json:"-"`
	EmailVerifySentAt *time.Time `gorm:"default:null" json:"-"`

	Password string `json:"-"`

	Sessions        []Session        `json:"sessions,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Profile         *Profile         `json:"profile,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PasswordReset   []PasswordReset  `json:"-" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SocialProviders []SocialProvider `json:"-" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return nil
}

// ComparePassword compares the password of a user
func (u *User) ComparePassword(password string) bool {
	return utils.CheckPasswordHash(password, u.Password)
}

// SetPassword sets the password of a user
func (u *User) SetPassword(password string) error {
	hashedPassword, err := utils.GeneratePasswordHash(password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword
	return nil
}

func (u *User) UserSessionLimiter() error {
	limit, err := strconv.Atoi(os.Getenv("MAX_SESSIONS_PER_USER"))

	if len(u.Sessions) >= limit {
		if err = DB.Delete(&u.Sessions[0]).Error; err != nil {
			return err
		}
	}
	return nil
}

// CheckPasswordResetToken checks the password reset code of a user
func (u *User) CheckPasswordResetToken(token string) bool {
	if len(u.PasswordReset) > 0 {
		return u.PasswordReset[0].CheckResetToken(token)
	}
	return false
}

// SetEmailVerifyCode verifies the email of a user
func (u *User) SetEmailVerifyCode(code string) {
	u.EmailVerifyCode = code
	currentTime := utils.GetCurrentTime()
	u.EmailVerifySentAt = &currentTime
}

// VerifyEmail verifies the email of a user
func (u *User) VerifyEmail(code string) bool {
	if u.EmailVerifyCode == code {
		u.EmailVerified = true
		u.EmailVerifyCode = ""
		u.EmailVerifySentAt = nil
		return true
	}
	return false
}

// UpdateUserEmail updates the email of a user and sets the email verified to false and email verify code to null
func (u *User) UpdateUserEmail(email string) {
	u.EmailVerified = false
	u.EmailVerifyCode = ""
	u.Email = email
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

func (u *User) CheckPasswordResetCode(code string) bool {
	if len(u.PasswordReset) > 0 {
		return u.PasswordReset[0].CheckResetToken(code)
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

func (u *User) GetLastPasswordReset() (*PasswordReset, error) {
	if len(u.PasswordReset) > 0 {
		return &u.PasswordReset[0], nil
	}
	return nil, errors.New("password reset not found")
}

func userPreload() *gorm.DB {
	return DB.Preload(clause.Associations)
}

func (u *User) CreateSession(ctx context.Context) (*token.DefaultToken, error) {
	sClient := ctx.Value(client.Key).(*client.Client)
	jwt, expiresIn, err := token.CreateJWT(u.ID)
	if err != nil {
		return nil, err
	}

	rToken, refreshTokenExpireAt := token.CreateRefreshToken(u.ID)

	if err := u.UserSessionLimiter(); err != nil {
		return nil, err
	}

	var session = Session{
		UserID:       u.ID,
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
	likeOperator := getLikeOperator(DB)

	var user = User{}

	query := DB.Model(&user).
		Joins("LEFT JOIN profiles ON profiles.user_id = users.id").
		Where(fmt.Sprintf("users.name %s ?", likeOperator), "%"+md.Search+"%").
		Or(fmt.Sprintf("users.email %s ?", likeOperator), "%"+md.Search+"%").
		Or(fmt.Sprintf("profiles.timezone %s ?", likeOperator), "%"+md.Search+"%")

	query = addDatabaseSpecificConditions(query, likeOperator, md.Search)

	query = query.
		Or(fmt.Sprintf("profiles.gender %s ?", likeOperator), "%"+md.Search+"%").
		Or(fmt.Sprintf("profiles.locale %s ?", likeOperator), "%"+md.Search+"%").
		Or(fmt.Sprintf("profiles.education %s ?", likeOperator), "%"+md.Search+"%").
		Count(&count).
		Order(user.ConvertToOrder(md)).
		Offset(int(md.Offset())).
		Limit(int(md.PageSize)).
		Find(&users)

	if err := query.Error; err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

// getLikeOperator determines the appropriate LIKE operator based on the database dialect.
func getLikeOperator(db *gorm.DB) string {
	if db.Dialector.Name() == "postgres" {
		return "ILIKE"
	}
	return "LIKE"
}

// addDatabaseSpecificConditions adds conditions to the query that are specific to the database dialect.
func addDatabaseSpecificConditions(query *gorm.DB, likeOperator, search string) *gorm.DB {
	if DB.Dialector.Name() == "postgres" {
		query = query.Or(fmt.Sprintf("profiles.birthdate::text %s ?", likeOperator), "%"+search+"%")
	} else {
		// Adjust for MySQL or other databases as necessary.
		query = query.Or(fmt.Sprintf("profiles.birthdate %s ?", likeOperator), "%"+search+"%")
	}
	return query
}
