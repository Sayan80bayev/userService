package service

import "userService/internal/models"

type ModerService interface {
	SetRoleById(userId uint, role models.Role) error

	BanUserById(userId uint) error

	UnBanUserById(userId uint) error
}
