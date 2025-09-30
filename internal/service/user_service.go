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
	"userService/internal/transport/request"
	"userService/internal/transport/response"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUserById(ctx context.Context, userId uuid.UUID) error
	GetAllUsers(ctx context.Context) ([]model.User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*model.User, error)
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

func (s *UserService) UpdateUser(ctx context.Context, ur request.UserRequest, userID uuid.UUID) error {
	u, err := s.userRepo.GetUserById(ctx, userID)
	if err != nil {
		return err
	}
	if u == nil {
		return fmt.Errorf("user not found: %s", userID)
	}

	oldURL := u.AvatarURL

	// Avatar update
	if ur.Avatar != nil && ur.Header != nil {
		if u.AvatarURL, err = s.fileStorage.UploadFile(ctx, ur.Avatar, ur.Header); err != nil {
			return err
		}
	}

	// Mandatory fields
	u.Lastname = ur.Lastname
	u.Firstname = ur.Firstname

	// Optional fields
	u.About = ur.About
	if ur.DateOfBirth == "" {
		u.DateOfBirth = nil
	} else {
		dob, err := date.ParseDate(ur.DateOfBirth)
		if err != nil {
			return err
		}
		u.DateOfBirth = &dob
	}
	u.Gender = ur.Gender
	u.Location = ur.Location
	u.Socials = ur.Socials

	// Persist update
	if err = s.userRepo.UpdateUser(ctx, u); err != nil {
		logging.Instance.Errorf("failed to update user %s: %v", userID, err)
		return err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("user:%s", userID)
	if err = s.cache.Delete(ctx, cacheKey); err != nil {
		logging.Instance.Warnf("failed to invalidate cache for user %s: %v", userID, err)
	}

	// Publish event (non-blocking for DB update)
	if err = s.producer.Produce(ctx, events.UserUpdated, events.UserUpdatedPayload{
		UserID:    userID,
		OldURL:    oldURL,
		AvatarURL: u.AvatarURL,
	}); err != nil {
		logging.Instance.Errorf("failed to publish UserUpdated event for user %s: %v", userID, err)
	}

	return nil
}

func (s *UserService) DeleteUserById(ctx context.Context, userId uuid.UUID) error {
	if err := s.userRepo.DeleteUserById(ctx, userId); err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("user:%s", userId)
	if err := s.cache.Delete(ctx, cacheKey); err != nil {
		logging.Instance.Warnf("failed to invalidate cache for deleted user %s: %v", userId, err)
	}

	// Publish event (non-blocking for DB delete)
	if err := s.producer.Produce(ctx, events.UserDeleted, events.UserDeletedPayload{
		UserID: userId,
	}); err != nil {
		logging.Instance.Errorf("failed to publish UserDeleted event for user %s: %v", userId, err)
	}

	return nil
}

func (s *UserService) GetUserById(ctx context.Context, id uuid.UUID) (*response.UserResponse, error) {
	cacheKey := fmt.Sprintf("user:%s", id.String())

	// 1. Try cache first
	if cached, err := s.cache.Get(ctx, cacheKey); err == nil && cached != "" {
		var ur response.UserResponse
		if err := json.Unmarshal([]byte(cached), &ur); err == nil {
			return &ur, nil
		}
		logging.Instance.Warnf("failed to unmarshal cached user %s: %v", id, err)
	} else if err != nil {
		logging.Instance.Warnf("cache get failed for user %s: %v", id, err)
	}

	// 2. Cache miss → get from DB
	user, err := s.userRepo.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	ur := s.mapper.Map(*user)

	// 3. Store in cache with TTL
	if data, err := json.Marshal(ur); err == nil {
		if err := s.cache.Set(ctx, cacheKey, data, 10*time.Minute); err != nil {
			logging.Instance.Warnf("failed to set cache for user %s: %v", id, err)
		}
	} else {
		logging.Instance.Warnf("failed to marshal user %s for cache: %v", id, err)
	}

	// 4. Return the user
	return &ur, nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]response.UserResponse, error) {
	users, err := s.userRepo.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	return s.mapper.MapEach(users), nil
}
