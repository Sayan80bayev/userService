package impl

import (
	"gorm.io/gorm"
	"userService/internal/models"
	"userService/internal/repository"
)

type ModerRepositoryImpl struct {
	db *gorm.DB
}

func NewModerRepoImpl(db *gorm.DB) repository.ModerRepository {
	return &ModerRepositoryImpl{db: db}
}

func (r *ModerRepositoryImpl) SetRoleById(userId uint, role models.Role) error {
	return r.db.Model(&models.User{}).Where("id = ?", userId).Update("role", role).Error
}

func (r *ModerRepositoryImpl) BanUserById(userId uint) error {
	return r.db.Model(&models.User{}).Where("id = ?", userId).Update("active", false).Error
}

func (r *ModerRepositoryImpl) UnBanUserById(userId uint) error {
	return r.db.Model(&models.User{}).Where("id = ?", userId).Update("active", true).Error
}
