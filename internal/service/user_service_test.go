package service

import (
	"bytes"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"mime/multipart"
	"testing"
	"time"
	"userService/internal/events"
	"userService/internal/model"
	"userService/internal/transfer/request"
)

type MockCacheService struct {
	mock.Mock
}

func (m *MockCacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockCacheService) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockCacheService) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockCacheService) Publish(ctx context.Context, channel, message string) error {
	args := m.Called(ctx, channel, message)
	return args.Error(0)
}

func (m *MockCacheService) Exists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func (m *MockCacheService) Subscribe(ctx context.Context, channel string) *redis.PubSub {
	args := m.Called(ctx, channel)
	if ps, ok := args.Get(0).(*redis.PubSub); ok {
		return ps
	}
	return nil
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUser(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteUserById(userId uuid.UUID) error {
	args := m.Called(userId)
	return args.Error(0)
}

func (m *MockUserRepository) GetAllUsers() ([]model.User, error) {
	args := m.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserRepository) GetUserById(id uuid.UUID) (*model.User, error) {
	args := m.Called(id)
	return args.Get(0).(*model.User), args.Error(1)
}

type MockFileService struct {
	mock.Mock
}

func (m *MockFileService) DeleteFileByURL(fileURL string) error {
	args := m.Called(fileURL)
	return args.Error(0)
}

func (m *MockFileService) UploadFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	args := m.Called(file, header)
	return args.String(0), args.Error(1)
}

type MockProducer struct {
	mock.Mock
}

func (m *MockProducer) Produce(topic string, event interface{}) error {
	args := m.Called(topic, event)
	return args.Error(0)
}

func (m *MockProducer) Close() {}

func TestUserService_UpdateUser(t *testing.T) {

	avatarFile, avatarHeader, err := createMockFile()
	if err != nil {
		t.Fatalf("Failed to create mock file: %v", err)
	}

	userUUID := uuid.New()
	tests := []struct {
		name          string
		setupMocks    func(*MockUserRepository, *MockFileService, *MockProducer)
		req           request.UserRequest
		userID        uuid.UUID
		expectedError string
	}{
		{
			name: "successful update",
			setupMocks: func(repo *MockUserRepository, fs *MockFileService, p *MockProducer) {
				repo.On("GetUserById", userUUID).Return(&model.User{AvatarURL: "old.jpg"}, nil)
				fs.On("UploadFile", mock.Anything, mock.Anything).Return("new.jpg", nil)
				repo.On("UpdateUser", mock.Anything).Return(nil)
				p.On("Produce", events.UserUpdated, mock.Anything).Return(nil)
			},
			req: request.UserRequest{
				Avatar:      avatarFile,
				Header:      avatarHeader,
				Firstname:   "newfirstname",
				Lastname:    "newlastname",
				DateOfBirth: "02.01.2004",
				About:       "about",
			},
			userID:        userUUID,
			expectedError: "",
		},
		{
			name: "error uploading file",
			setupMocks: func(repo *MockUserRepository, fs *MockFileService, p *MockProducer) {
				repo.On("GetUserById", userUUID).Return(&model.User{}, nil)
				fs.On("UploadFile", mock.Anything, mock.Anything).Return("", errors.New("upload error"))
			},
			req: request.UserRequest{
				Avatar: avatarFile,
				Header: avatarHeader,
			},
			userID:        userUUID,
			expectedError: "upload error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockUserRepository)
			fs := new(MockFileService)
			p := new(MockProducer)
			tt.setupMocks(repo, fs, p)

			svc := NewUserService(repo, fs, p, nil)
			err := svc.UpdateUser(tt.req, tt.userID)

			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError)
			}
		})
	}
}

func TestUserService_DeleteUserById(t *testing.T) {
	userUUID := uuid.New()
	repo := new(MockUserRepository)
	repo.On("DeleteUserById", userUUID).Return(nil)

	svc := NewUserService(repo, nil, nil, nil)
	err := svc.DeleteUserById(userUUID)

	assert.NoError(t, err)
}

func TestUserService_GetUserById(t *testing.T) {
	cache := new(MockCacheService)

	cache.On("Get", mock.Anything, mock.Anything).Return("", errors.New("cache miss"))
	cache.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	userUUID := uuid.New()
	repo := new(MockUserRepository)
	user := &model.User{
		ID:        userUUID,
		Firstname: "testuser",
	}
	repo.On("GetUserById", userUUID).Return(user, nil)

	svc := NewUserService(repo, nil, nil, cache)
	resp, err := svc.GetUserById(context.Background(), user.ID)

	assert.NoError(t, err)
	assert.Equal(t, user.ID, resp.ID)
	assert.Equal(t, "testuser", resp.Firstname)
}

func TestUserService_GetAllUsers(t *testing.T) {
	userUUID1 := uuid.New()
	userUUID2 := uuid.New()
	repo := new(MockUserRepository)
	users := []model.User{
		{
			ID:        userUUID1,
			Firstname: "user1",
		},
		{
			ID:        userUUID2,
			Firstname: "user2",
		},
	}
	repo.On("GetAllUsers").Return(users, nil)

	svc := NewUserService(repo, nil, nil, nil)
	resp, err := svc.GetAllUsers()

	assert.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, userUUID1, resp[0].ID)
	assert.Equal(t, "user1", resp[0].Firstname)
	assert.Equal(t, userUUID2, resp[1].ID)
	assert.Equal(t, "user2", resp[1].Firstname)
}

type mockMultipartFile struct {
	*bytes.Reader
}

func (m *mockMultipartFile) Close() error {
	return nil
}

// Создаем мок для файла и его заголовка.
func createMockFile() (multipart.File, *multipart.FileHeader, error) {
	// Создаем буфер, который будет имитировать файл
	fileContent := []byte("This is a mock file content")
	fileReader := bytes.NewReader(fileContent)

	// Создаем mockMultipartFile, который реализует multipart.File
	mockFile := &mockMultipartFile{Reader: fileReader}

	// Заголовок файла
	fileHeader := &multipart.FileHeader{
		Filename: "mockfile.txt",
		Size:     int64(len(fileContent)),
	}

	return mockFile, fileHeader, nil
}
