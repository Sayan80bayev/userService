package service

import (
	"userService/internal/model"
)

type ModerRepository interface {
	SetRoleById(userId int, role model.Role) error

	BanUserById(userId int) error

	UnBanUserById(userId int) error
}

type ModerService struct {
	repo ModerRepository
}

func NewModerService(moderRepository ModerRepository) *ModerService {
	return &ModerService{repo: moderRepository}
}

func (m *ModerService) SetRoleById(userId int, role model.Role) error {
	return m.repo.SetRoleById(userId, role)
}

func (m *ModerService) BanUserById(userId int) error {
	return m.repo.BanUserById(userId)
}

func (m *ModerService) UnBanUserById(userId int) error {
	return m.repo.UnBanUserById(userId)
}
