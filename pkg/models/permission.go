package UserModel

import (
	"acourse-auth-user-service/pkg/database"
	"time"
)

type Permission struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	Name      string     `json:"name"`
	Roles     []*Role    `gorm:"many2many:role_permissions;"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

func init() {
	database.Connect()
	dbConnection := database.GetConnection()
	dbConnection.AutoMigrate(&Permission{})
}
