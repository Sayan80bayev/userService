package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username    string `gorm:"unique"`
	Password    string
	About       string
	Active      bool `gorm:"default:true"`
	DateOfBirth time.Time
	AvatarURL   string
}
