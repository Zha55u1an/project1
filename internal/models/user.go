package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	// ID       int       `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
}
