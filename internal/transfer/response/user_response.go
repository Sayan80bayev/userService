package response

import (
	"gorm.io/gorm"
	"time"
)

type UserResponse struct {
	gorm.Model
	Username    string    `json:"username"`
	About       string    `json:"about,omitempty"`
	Active      bool      `json:"active"`
	DateOfBirth time.Time `json:"date,omitempty"`
	AvatarURL   string    `json:"avatar_url,omitempty"`
}
