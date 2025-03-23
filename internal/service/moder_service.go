package service

import (
	"errors"
	"userService/internal/model"
)

type ModerRepository interface {
	GetRoleById(userId int) (model.Role, error)
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

// Проверяет, может ли редактор управлять указанной ролью
func (m *ModerService) validateModeration(editorRole model.Role, userId int) error {
	targetRole, err := m.repo.GetRoleById(userId)
	if err != nil {
		return errors.New("failed to fetch target user role")
	}

	if !editorRole.CanModerate(targetRole) {
		return errors.New("you do not have permission to perform this action")
	}

	return nil
}

// SetRoleById меняет роль пользователя, учитывая иерархию.
func (m *ModerService) SetRoleById(editorId int, editorRole model.Role, userId int, newRole model.Role) error {
	if err := m.validateModeration(editorRole, userId); err != nil {
		return err
	}
	if !editorRole.CanModerate(newRole) {
		return errors.New("you do not have permission to perform this action")
	}

	return m.repo.SetRoleById(userId, newRole)
}

// BanUserById блокирует пользователя, проверяя иерархию ролей.
func (m *ModerService) BanUserById(editorRole model.Role, userId int) error {
	if err := m.validateModeration(editorRole, userId); err != nil {
		return err
	}
	return m.repo.BanUserById(userId)
}

// UnBanUserById снимает блокировку пользователя, проверяя иерархию ролей.
func (m *ModerService) UnBanUserById(editorRole model.Role, userId int) error {
	if err := m.validateModeration(editorRole, userId); err != nil {
		return err
	}
	return m.repo.UnBanUserById(userId)
}
