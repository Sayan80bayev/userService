package service

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"userService/internal/events"
	"userService/internal/mappers"
	"userService/internal/messaging"
	"userService/internal/model"
	"userService/internal/pkg/storage"
	"userService/internal/transfer/request"
	"userService/internal/transfer/response"
	"userService/pkg/date"
	"userService/pkg/logging"
	"userService/pkg/mapping"
)

type UserRepository interface {
	CreateUser(user *model.User) error

	UpdateUser(user *model.User) error

	DeleteUserById(userId int) error

	GetAllUsers() ([]model.User, error)

	GetUserByUsername(username string) (*model.User, error)

	GetUserById(id int) (*model.User, error)
}

type UserService struct {
	userRepo    UserRepository
	fileService storage.FileService
	producer    messaging.Producer
	mapper      mapping.MapFunc[model.User, response.UserResponse]
}

func NewUserService(
	userRepo UserRepository,
	fileService storage.FileService,
	producer messaging.Producer,
) *UserService {
	return &UserService{
		userRepo:    userRepo,
		fileService: fileService,
		producer:    producer,
		mapper:      mappers.UserToUserResponse,
	}
}

func (s *UserService) UpdateUser(ur request.UserRequest, userID int) error {
	u, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return err
	}

	oldURL := u.AvatarURL
	if ur.Avatar != nil && ur.Header != nil {
		if u.AvatarURL, err = s.fileService.UploadFile(ur.Avatar, ur.Header); err != nil {
			return err
		}
	}

	dob, err := date.ParseDate(ur.DateOfBirth)
	if err != nil {
		return err
	}

	u.Username, u.DateOfBirth, u.About = ur.Username, dob, ur.About

	err = s.userRepo.UpdateUser(u)
	if err != nil {
		logging.Instance.Error(err)
		return err
	}

	return s.producer.Produce("UserUpdate", events.UserUpdated{
		UserID:    userID,
		OldURL:    oldURL,
		AvatarURL: u.AvatarURL,
	})
}

func (s *UserService) DeleteUserById(userId int) error {
	return s.userRepo.DeleteUserById(userId)
}

func (s *UserService) GetUserByUsername(username string) (*response.UserResponse, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	ur := s.mapper.Map(*user)
	return &ur, nil
}

func (s *UserService) ChangePassword(userId int, pr request.ChangePasswordRequest) error {
	user, err := s.userRepo.GetUserById(userId)
	if err != nil {
		logging.Instance.Warn(err)
		return errors.New("Could not find user. ")
	}

	if pr.ConfirmPassword != pr.NewPassword {
		return errors.New("Passwords are not the same. ")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pr.OldPassword)); err != nil {
		return errors.New("Invalid credentials. ")
	}

	user.Password = pr.NewPassword
	return s.userRepo.UpdateUser(user)
}

func (s *UserService) GetUserById(id int) (*response.UserResponse, error) {
	user, err := s.userRepo.GetUserById(id)
	if err != nil {
		return nil, err
	}
	ur := s.mapper.Map(*user)
	return &ur, nil
}

func (s *UserService) GetAllUsers() ([]response.UserResponse, error) {
	users, err := s.userRepo.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return s.mapper.MapEach(users), nil
}
