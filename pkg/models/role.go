package model

import (
	"acourse-auth-user-service/pkg/database"
	"time"
)

type Role struct {
	ID          uint         `gorm:"primary_key" json:"id,omitempty"`
	Name        string       `gorm:"unique" json:"name,omitempty"`
	Users       []User       `gorm:"foreignKey:RoleID" json:"users,omitempty"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	CreatedAt   time.Time    `json:"created_at,omitempty"`
	UpdatedAt   time.Time    `json:"updated_at,omitempty"`
	DeletedAt   *time.Time   `sql:"index" json:"deleted_at,omitempty"`
}

func FindRole(id uint) (Role, error) {

	var role Role
	err := database.GetConnection().First(&role, id).Error

	if err != nil {
		return role, err
	}

	return role, nil

}
