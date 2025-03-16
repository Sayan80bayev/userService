package service

import "userService/internal/models"

type UserService interface {
	UpdateUser(user *models.User) error
	DeleteUserById(userId uint) error

	GetUserByUsername(username string) (*models.User, error)
}
