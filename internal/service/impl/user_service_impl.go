package impl

import (
	"userService/internal/dto/request"
	"userService/internal/models"
	"userService/internal/repository"
	"userService/internal/service"
)

type UserServiceImpl struct {
	userRepo repository.UserRepository
}

func (impl *UserServiceImpl) UpdateUser(ur *request.UserRequest) error {
	//TODO implement me
	panic("implement me")
}

func NewUserServiceImpl(userRepo repository.UserRepository) service.UserService {
	return &UserServiceImpl{userRepo: userRepo}
}

func (impl *UserServiceImpl) DeleteUserById(userId uint) error {
	//TODO implement me
	panic("implement me")
}

func (impl *UserServiceImpl) GetUserByUsername(username string) (*models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (impl *UserServiceImpl) GetAllUsers() ([]*models.User, error) {
	panic("implement me")
}
