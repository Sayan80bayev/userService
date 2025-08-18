package service

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"mime/multipart"
	"testing"
	"userService/internal/model"
	"userService/internal/transfer/request"
)

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

func (m *MockUserRepository) DeleteUserById(userId int) error {
	args := m.Called(userId)
	return args.Error(0)
}

func (m *MockUserRepository) GetAllUsers() ([]model.User, error) {
	args := m.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByUsername(username string) (*model.User, error) {
	args := m.Called(username)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetUserById(id int) (*model.User, error) {
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

	tests := []struct {
		name          string
		setupMocks    func(*MockUserRepository, *MockFileService, *MockProducer)
		req           request.UserRequest
		userID        int
		expectedError string
	}{
		{
			name: "successful update",
			setupMocks: func(repo *MockUserRepository, fs *MockFileService, p *MockProducer) {
				repo.On("GetUserById", 1).Return(&model.User{AvatarURL: "old.jpg"}, nil)
				// Обратите внимание, что мок теперь принимает правильные типы
				fs.On("UploadFile", mock.Anything, mock.Anything).Return("new.jpg", nil)
				repo.On("UpdateUser", mock.Anything).Return(nil)
				p.On("Produce", "UserUpdate", mock.Anything).Return(nil)
			},
			req: request.UserRequest{
				Avatar:      avatarFile,
				Header:      avatarHeader,
				Username:    "newuser",
				DateOfBirth: "02.01.2004",
				About:       "about",
			},
			userID:        1,
			expectedError: "",
		},
		{
			name: "error uploading file",
			setupMocks: func(repo *MockUserRepository, fs *MockFileService, p *MockProducer) {
				repo.On("GetUserById", 1).Return(&model.User{}, nil)
				fs.On("UploadFile", mock.Anything, mock.Anything).Return("", errors.New("upload error"))
			},
			req: request.UserRequest{
				Avatar: avatarFile,
				Header: avatarHeader,
			},
			userID:        1,
			expectedError: "upload error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockUserRepository)
			fs := new(MockFileService)
			p := new(MockProducer)
			tt.setupMocks(repo, fs, p)

			svc := NewUserService(repo, fs, p)
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
	repo := new(MockUserRepository)
	repo.On("DeleteUserById", 1).Return(nil)

	svc := NewUserService(repo, nil, nil)
	err := svc.DeleteUserById(1)

	assert.NoError(t, err)
}

func TestUserService_GetUserByUsername(t *testing.T) {
	repo := new(MockUserRepository)
	user := &model.User{Username: "test"}
	repo.On("GetUserByUsername", "test").Return(user, nil)

	svc := NewUserService(repo, nil, nil)
	resp, err := svc.GetUserByUsername("test")

	assert.NoError(t, err)
	assert.Equal(t, "test", resp.Username)
}

func TestUserService_GetUserById(t *testing.T) {
	repo := new(MockUserRepository)
	user := &model.User{
		Model:    gorm.Model{ID: 1},
		Username: "testuser",
	}
	repo.On("GetUserById", 1).Return(user, nil)

	svc := NewUserService(repo, nil, nil)
	resp, err := svc.GetUserById(1)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), resp.ID)
	assert.Equal(t, "testuser", resp.Username)
}

func TestUserService_GetAllUsers(t *testing.T) {
	repo := new(MockUserRepository)
	users := []model.User{
		{
			Model:    gorm.Model{ID: 1},
			Username: "user1",
		},
		{
			Model:    gorm.Model{ID: 2},
			Username: "user2",
		},
	}
	repo.On("GetAllUsers").Return(users, nil)

	svc := NewUserService(repo, nil, nil)
	resp, err := svc.GetAllUsers()

	assert.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, uint(1), resp[0].ID)
	assert.Equal(t, "user1", resp[0].Username)
	assert.Equal(t, uint(2), resp[1].ID)
	assert.Equal(t, "user2", resp[1].Username)
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
