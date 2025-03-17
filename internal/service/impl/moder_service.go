package impl

import (
	"userService/internal/models"
	"userService/internal/repository"
	"userService/internal/service"
)

type ModerServiceImpl struct {
	userRepo repository.UserRepository
}

func NewModerServiceImpl(userRepo repository.UserRepository) service.ModerService {
	return &ModerServiceImpl{}
}

func (m ModerServiceImpl) SetRoleById(userId uint, role models.Role) error {
	//TODO implement me
	panic("implement me")
}

func (m ModerServiceImpl) BanUserById(userId uint) error {
	//TODO implement me
	panic("implement me")
}

func (m ModerServiceImpl) UnBanUserById(userId uint) error {
	//TODO implement me
	panic("implement me")
}
