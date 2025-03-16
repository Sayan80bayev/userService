package repository

import (
	"userService/internal/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error

	UpdateUser(user *models.User) error

	DeleteUserById(userId uint) error

	GetUserByUsername(username string) (*models.User, error)

	SetRoleById(userId uint, role models.Role) error

	BanUserById(userId uint) error

	UnBanUserById(userId uint) error
}
