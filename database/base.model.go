package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type UUIDBaseModel struct {
	ID        uuid.UUID       `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type UINTBaseModel struct {
	ID        uint64          `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
