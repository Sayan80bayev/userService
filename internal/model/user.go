package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model  `swaggerignore:"true"`
	Email       string `gorm:"unique"`
	Username    string `gorm:"unique"`
	About       string
	Active      bool `gorm:"default:true"`
	DateOfBirth time.Time
	AvatarURL   string
}
