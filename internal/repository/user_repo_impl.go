package repository

import (
	"sync"

	"gorm.io/gorm"
	"userService/internal/model"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

var (
	instance *UserRepositoryImpl
	once     sync.Once
)

// NewUserRepository returns a singleton instance of UserRepositoryImpl
func NewUserRepository(db *gorm.DB) *UserRepositoryImpl {
	once.Do(func() {
		instance = &UserRepositoryImpl{db: db}
	})
	return instance
}

func (r *UserRepositoryImpl) GetAllUsers() ([]model.User, error) {
	var users []model.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *UserRepositoryImpl) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepositoryImpl) UpdateUser(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepositoryImpl) DeleteUserById(userId int) error {
	return r.db.Delete(&model.User{}, userId).Error
}

func (r *UserRepositoryImpl) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *UserRepositoryImpl) GetUserById(id int) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", id).First(&user).Error
	return &user, err
}
