package database

import (
	"github.com/google/uuid"
	"github.com/usercoredev/usercore/utils"
	"gorm.io/gorm"
	"log"
	"reflect"
	"time"
)

type BaseModel struct {
	OrderBy string
	Order   string
}

type UUIDBaseModel struct {
	BaseModel
	ID        uuid.UUID       `gorm:"primaryKey" json:"id" sortable:"true"`
	CreatedAt time.Time       `json:"created_at" sortable:"true"`
	UpdatedAt time.Time       `json:"updated_at" sortable:"true"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type UINTBaseModel struct {
	BaseModel
	ID        uint64          `gorm:"primaryKey" json:"id" sortable:"true"`
	CreatedAt time.Time       `json:"created_at" sortable:"true"`
	UpdatedAt time.Time       `json:"updated_at" sortable:"true"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func getSortableFields(v interface{}) map[string]bool {
	fields := make(map[string]bool)
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		if sortable, ok := field.Tag.Lookup("sortable"); ok && sortable == "true" {
			fields[field.Name] = true
		}
	}
	return fields
}

func (bm *BaseModel) ConvertToOrder(metadata utils.PageMetadata) string {
	if _, ok := getSortableFields(bm)[metadata.OrderBy]; !ok {
		log.Println("Invalid order by field, defaulting to created_at")
		metadata.OrderBy = "created_at"
	}

	safeOrderValues := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	if _, ok := safeOrderValues[bm.Order]; !ok {
		metadata.Order = "desc"
	}

	return metadata.OrderBy + " " + metadata.Order
}
