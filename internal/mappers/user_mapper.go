package mappers

import (
	"github.com/Sayan80bayev/go-project/pkg/mapper"
	"userService/internal/model"
	"userService/internal/transfer/response"
)

type UserMapper struct {
	mapper.MapFunc[model.User, response.UserResponse]
}

func NewUserMapper() *UserMapper {
	return &UserMapper{MapFunc: UserToUserResponse}
}

// UserToUserResponse maps a User to UserResponse
var UserToUserResponse = mapper.MapFunc[model.User, response.UserResponse](func(u model.User) response.UserResponse {
	return response.UserResponse{
		ID:          u.ID, // Copies ID, CreatedAt, UpdatedAt, DeletedAt from gorm.Model
		Username:    u.Username,
		About:       u.About,
		DateOfBirth: u.DateOfBirth,
		AvatarURL:   u.AvatarURL,
		Email:       u.Email,
		Socials:     u.Socials,
		Gender:      u.Gender,
		Location:    u.Location,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		DeletedAt:   u.DeletedAt,
	}
})
