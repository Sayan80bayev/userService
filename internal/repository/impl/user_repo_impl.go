package impl

import (
	"gorm.io/gorm"
	"userService/internal/models"
	"userService/internal/repository"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepositoryImpl{db}
}

func (r *UserRepositoryImpl) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepositoryImpl) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepositoryImpl) DeleteUserById(userId uint) error {
	return r.db.Delete(&models.User{}, userId).Error
}

func (r *UserRepositoryImpl) SetRoleById(userId uint, role models.Role) error {
	return r.db.Model(&models.User{}).Where("id = ?", userId).Update("role", role).Error
}

func (r *UserRepositoryImpl) BanUserById(userId uint) error {
	return r.db.Model(&models.User{}).Where("id = ?", userId).Update("active", false).Error
}

func (r *UserRepositoryImpl) UnBanUserById(userId uint) error {
	return r.db.Model(&models.User{}).Where("id = ?", userId).Update("active", true).Error
}

func (r *UserRepositoryImpl) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}
