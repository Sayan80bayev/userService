package service

import (
	"userService/internal/model"
)

type ModerRepository interface {
	SetRoleById(userId uint, role model.Role) error

	BanUserById(userId uint) error

	UnBanUserById(userId uint) error
}

type ModerService struct {
	moderRepository ModerRepository
}

func NewModerServiceImpl(moderRepository ModerRepository) *ModerService {
	return &ModerService{moderRepository: moderRepository}
}

func (m *ModerService) SetRoleById(userId uint, role model.Role) error {
	//TODO implement me
	panic("implement me")
}

func (m *ModerService) BanUserById(userId uint) error {
	//TODO implement me
	panic("implement me")
}

func (m *ModerService) UnBanUserById(userId uint) error {
	//TODO implement me
	panic("implement me")
}
