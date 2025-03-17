package service

import (
	"userService/internal/dto/request"
	"userService/internal/models"
)

type UserService interface {
	UpdateUser(ur *request.UserRequest) error

	DeleteUserById(userId uint) error

	GetUserByUsername(username string) (*models.User, error)
	
	GetAllUsers() ([]*models.User, error)
}
