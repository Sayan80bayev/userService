package mappers

import (
	"userService/internal/model"
	"userService/internal/response"
	"userService/pkg/mapping"
)

type UserMapper struct {
	mapping.MapFunc[model.User, response.UserResponse]
}

// UserToUserResponse maps a User to UserResponse
var UserToUserResponse = mapping.MapFunc[model.User, response.UserResponse](func(u model.User) response.UserResponse {
	return response.UserResponse{
		Model:       u.Model, // Copies ID, CreatedAt, UpdatedAt, DeletedAt from gorm.Model
		Username:    u.Username,
		About:       u.About,
		Active:      u.Active,
		DateOfBirth: u.DateOfBirth,
		AvatarURL:   u.AvatarURL,
	}
})
