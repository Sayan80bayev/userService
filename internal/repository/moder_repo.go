package repository

import "userService/internal/models"

type ModerRepository interface {
	SetRoleById(userId uint, role models.Role) error

	BanUserById(userId uint) error

	UnBanUserById(userId uint) error
}
