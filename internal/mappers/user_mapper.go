package mappers

import (
	"github.com/Sayan80bayev/go-project/pkg/mapper"
	"userService/internal/model"
	"userService/internal/transfer/response"
)

type UserMapper struct {
	mapper.MapFunc[model.User, response.UserResponse]
}

// UserToUserResponse maps a User to UserResponse
var UserToUserResponse = mapper.MapFunc[model.User, response.UserResponse](func(u model.User) response.UserResponse {
	return response.UserResponse{
		Model:       u.Model, // Copies ID, CreatedAt, UpdatedAt, DeletedAt from gorm.Model
		Username:    u.Username,
		About:       u.About,
		Active:      u.Active,
		DateOfBirth: u.DateOfBirth,
		AvatarURL:   u.AvatarURL,
	}
})
