package model

import (
	"time"
)

type Permission struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	Name      string     `json:"name"`
	Code      string     `json:"code"`
	Roles     []Role     `gorm:"many2many:role_permissions;" json:"roles,omitempty"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`
}

func (p *Permission) ConstructPermissionCode() {

}
