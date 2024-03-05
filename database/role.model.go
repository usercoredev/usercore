package database

type Role struct {
	UINTBaseModel
	Name        string       `json:"name" gorm:"type:varchar(255);not null"`
	Key         string       `json:"key" gorm:"type:varchar(255);not null;unique"`
	Description string       `json:"description" gorm:"type:varchar(255);not null"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
}
