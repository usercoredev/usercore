package database

import (
	"github.com/google/uuid"
	"time"
)

type Profile struct {
	UINTBaseModel
	UserID    uuid.UUID  `gorm:"default:null" json:"-"`
	Picture   string     `gorm:"default:null" json:"picture,omitempty"`
	Gender    string     `gorm:"default:null" json:"gender,omitempty"`
	Education string     `gorm:"default:null" json:"education,omitempty"`
	Birthdate *time.Time `gorm:"default:null" json:"birthdate,omitempty"`
	Locale    string     `gorm:"default:null" json:"locale,omitempty"`
	Timezone  string     `gorm:"default:null" json:"timezone,omitempty"`
}
