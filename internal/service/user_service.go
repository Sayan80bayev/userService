package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Sayan80bayev/go-project/pkg/caching"
	"github.com/Sayan80bayev/go-project/pkg/date"
	"github.com/Sayan80bayev/go-project/pkg/logging"
	"github.com/Sayan80bayev/go-project/pkg/messaging"
	storage "github.com/Sayan80bayev/go-project/pkg/objectStorage"
	"github.com/google/uuid"
	"time"
	"userService/internal/events"
	"userService/internal/mappers"
	"userService/internal/model"
	"userService/internal/transfer/request"
	"userService/internal/transfer/response"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	UpdateUser(user *model.User) error
	DeleteUserById(userId uuid.UUID) error
	GetAllUsers() ([]model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	GetUserById(id uuid.UUID) (*model.User, error)
}

type UserService struct {
	cache       caching.CacheService
	userRepo    UserRepository
	fileStorage storage.FileStorage
	producer    messaging.Producer
	mapper      *mappers.UserMapper
}

func NewUserService(
	userRepo UserRepository,
	fileStorage storage.FileStorage,
	producer messaging.Producer,
	cache caching.CacheService,
) *UserService {
	return &UserService{
		userRepo:    userRepo,
		fileStorage: fileStorage,
		producer:    producer,
		mapper:      mappers.NewUserMapper(),
		cache:       cache,
	}
}

func (s *UserService) UpdateUser(ur request.UserRequest, userID uuid.UUID) error {
	u, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return err
	}

	oldURL := u.AvatarURL
	if ur.Avatar != nil && ur.Header != nil {
		if u.AvatarURL, err = s.fileStorage.UploadFile(ur.Avatar, ur.Header); err != nil {
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

	return s.producer.Produce("UserUpdate", events.UserUpdatedPayload{
		UserID:    userID,
		OldURL:    oldURL,
		AvatarURL: u.AvatarURL,
	})
}

func (s *UserService) DeleteUserById(userId uuid.UUID) error {
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

func (s *UserService) GetUserById(ctx context.Context, id uuid.UUID) (*response.UserResponse, error) {
	cacheKey := fmt.Sprintf("user:%s", id.String())

	// 1. Try cache first
	if cached, err := s.cache.Get(ctx, cacheKey); err == nil && cached != "" {
		var ur response.UserResponse
		if err := json.Unmarshal([]byte(cached), &ur); err == nil {
			return &ur, nil
		}
		// if unmarshal fails, fall through to DB
	}

	// 2. Cache miss → get from DB
	user, err := s.userRepo.GetUserById(id)
	if err != nil {
		return nil, err
	}

	ur := s.mapper.Map(*user)

	// 3. Store in cache with TTL (e.g. 10 minutes)
	if data, err := json.Marshal(ur); err == nil {
		_ = s.cache.Set(ctx, cacheKey, data, 10*time.Minute)
	}

	// 4. Return the user
	return &ur, nil
}

func (s *UserService) GetAllUsers() ([]response.UserResponse, error) {
	users, err := s.userRepo.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return s.mapper.MapEach(users), nil
}
