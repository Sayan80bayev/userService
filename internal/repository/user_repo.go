package repository

import (
	"userService/internal/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error

	UpdateUser(user *models.User) error

	DeleteUserById(userId uint) error

	GetAllUsers() ([]*models.User, error)

	GetUserByUsername(username string) (*models.User, error)
}
