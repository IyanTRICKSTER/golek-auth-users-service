package UserModel

import (
	"acourse-auth-user-service/pkg/database"
	"time"
)

type Role struct {
	ID          uint          `gorm:"primary_key" json:"id,omitempty"`
	Name        string        `json:"name,omitempty"`
	Users       []User        `gorm:"foreignKey:ID" json:"users,omitempty"`
	Permissions []*Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	CreatedAt   time.Time     `json:"created_at,omitempty"`
	UpdatedAt   time.Time     `json:"updated_at,omitempty"`
	DeletedAt   *time.Time    `sql:"index" json:"deleted_at,omitempty"`
}

func init() {
	database.Connect()
	dbConnection := database.GetConnection()
	dbConnection.AutoMigrate(&Role{})
}
