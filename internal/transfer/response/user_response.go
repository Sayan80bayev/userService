package response

import (
	"github.com/google/uuid"
	"time"
)

type UserResponse struct {
	ID        uuid.UUID  `bson:"_id,omitempty" json:"id" validate:"omitempty"`
	CreatedAt time.Time  `bson:"created_at,omitempty" json:"created_at" validate:"omitempty"`
	UpdatedAt time.Time  `bson:"updated_at,omitempty" json:"updated_at" validate:"omitempty"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty" validate:"omitempty"`

	Email    string `bson:"email" json:"email" validate:"required,email"`
	Username string `bson:"username" json:"username" validate:"required,min=3,max=20"`
	About    string `bson:"about,omitempty" json:"about,omitempty" validate:"omitempty,max=500"`

	DateOfBirth *time.Time `bson:"date_of_birth,omitempty" json:"date_of_birth,omitempty" validate:"omitempty,lte"`
	AvatarURL   string     `bson:"avatar_url,omitempty" json:"avatar_url,omitempty" validate:"omitempty,url"`

	Gender   string `bson:"gender,omitempty" json:"gender,omitempty" validate:"omitempty,oneof=male female other"`
	Location string `bson:"location,omitempty" json:"location,omitempty" validate:"omitempty,max=100"`

	Socials []string `bson:"socials,omitempty" json:"socials,omitempty" validate:"omitempty,dive,url"`
}
