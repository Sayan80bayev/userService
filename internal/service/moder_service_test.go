package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"userService/internal/model"
)

// MockRepository реализует ModerRepository для тестов
type MockRepository struct {
	GetRoleByIdFunc   func(userId int) (model.Role, error)
	SetRoleByIdFunc   func(userId int, role model.Role) error
	BanUserByIdFunc   func(userId int) error
	UnBanUserByIdFunc func(userId int) error
}

func (m *MockRepository) GetRoleById(userId int) (model.Role, error) {
	return m.GetRoleByIdFunc(userId)
}

func (m *MockRepository) SetRoleById(userId int, role model.Role) error {
	return m.SetRoleByIdFunc(userId, role)
}

func (m *MockRepository) BanUserById(userId int) error {
	return m.BanUserByIdFunc(userId)
}

func (m *MockRepository) UnBanUserById(userId int) error {
	return m.UnBanUserByIdFunc(userId)
}

func TestModerService_SetRoleById(t *testing.T) {
	tests := []struct {
		name          string
		editorRole    model.Role
		targetRole    model.Role
		newRole       model.Role
		repoError     error
		expectedError string
	}{
		{
			name:          "Admin can set moderator role",
			editorRole:    model.RoleAdmin,
			targetRole:    model.RoleUser,
			newRole:       model.RoleModerator,
			expectedError: "",
		},
		{
			name:          "Moderator cannot set admin role",
			editorRole:    model.RoleModerator,
			targetRole:    model.RoleUser,
			newRole:       model.RoleAdmin,
			expectedError: "you do not have permission to perform this action",
		},
		{
			name:          "Moderator cannot moderate admin",
			editorRole:    model.RoleModerator,
			targetRole:    model.RoleAdmin,
			newRole:       model.RoleUser,
			expectedError: "you do not have permission to perform this action",
		},
		{
			name:          "User cannot moderate anyone",
			editorRole:    model.RoleUser,
			targetRole:    model.RoleUser,
			newRole:       model.RoleModerator,
			expectedError: "you do not have permission to perform this action",
		},
		{
			name:          "Error getting target role",
			editorRole:    model.RoleAdmin,
			targetRole:    model.RoleUser,
			newRole:       model.RoleModerator,
			repoError:     errors.New("database error"),
			expectedError: "failed to fetch target user role",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockRepository{
				GetRoleByIdFunc: func(userId int) (model.Role, error) {
					return tt.targetRole, tt.repoError
				},
				SetRoleByIdFunc: func(userId int, role model.Role) error {
					if tt.expectedError != "" {
						return errors.New(tt.expectedError)
					}
					return nil
				},
			}

			moderService := NewModerService(repo)
			err := moderService.SetRoleById(1, tt.editorRole, 2, tt.newRole)

			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError)
			}
		})
	}
}

func TestModerService_BanUserById(t *testing.T) {
	tests := []struct {
		name          string
		editorRole    model.Role
		targetRole    model.Role
		repoError     error
		expectedError string
	}{
		{
			name:          "Admin can ban user",
			editorRole:    model.RoleAdmin,
			targetRole:    model.RoleUser,
			expectedError: "",
		},
		{
			name:          "Moderator can ban user",
			editorRole:    model.RoleModerator,
			targetRole:    model.RoleUser,
			expectedError: "",
		},
		{
			name:          "Moderator cannot ban admin",
			editorRole:    model.RoleModerator,
			targetRole:    model.RoleAdmin,
			expectedError: "you do not have permission to perform this action",
		},
		{
			name:          "User cannot ban anyone",
			editorRole:    model.RoleUser,
			targetRole:    model.RoleUser,
			expectedError: "you do not have permission to perform this action",
		},
		{
			name:          "Error getting target role",
			editorRole:    model.RoleAdmin,
			targetRole:    model.RoleUser,
			repoError:     errors.New("database error"),
			expectedError: "failed to fetch target user role",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockRepository{
				GetRoleByIdFunc: func(userId int) (model.Role, error) {
					return tt.targetRole, tt.repoError
				},
				BanUserByIdFunc: func(userId int) error {
					return nil
				},
			}

			moderService := NewModerService(repo)
			err := moderService.BanUserById(tt.editorRole, 2)

			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError)
			}
		})
	}
}

func TestModerService_UnBanUserById(t *testing.T) {
	tests := []struct {
		name          string
		editorRole    model.Role
		targetRole    model.Role
		repoError     error
		expectedError string
	}{
		{
			name:          "Admin can unban user",
			editorRole:    model.RoleAdmin,
			targetRole:    model.RoleUser,
			expectedError: "",
		},
		{
			name:          "Moderator can unban user",
			editorRole:    model.RoleModerator,
			targetRole:    model.RoleUser,
			expectedError: "",
		},
		{
			name:          "Moderator cannot unban admin",
			editorRole:    model.RoleModerator,
			targetRole:    model.RoleAdmin,
			expectedError: "you do not have permission to perform this action",
		},
		{
			name:          "User cannot unban anyone",
			editorRole:    model.RoleUser,
			targetRole:    model.RoleUser,
			expectedError: "you do not have permission to perform this action",
		},
		{
			name:          "Error getting target role",
			editorRole:    model.RoleAdmin,
			targetRole:    model.RoleUser,
			repoError:     errors.New("database error"),
			expectedError: "failed to fetch target user role",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockRepository{
				GetRoleByIdFunc: func(userId int) (model.Role, error) {
					return tt.targetRole, tt.repoError
				},
				UnBanUserByIdFunc: func(userId int) error {
					return nil
				},
			}

			moderService := NewModerService(repo)
			err := moderService.UnBanUserById(tt.editorRole, 2)

			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError)
			}
		})
	}
}

func TestValidateModeration(t *testing.T) {
	tests := []struct {
		name          string
		editorRole    model.Role
		targetRole    model.Role
		repoError     error
		expectedError string
	}{
		{
			name:          "Admin can moderate moderator",
			editorRole:    model.RoleAdmin,
			targetRole:    model.RoleModerator,
			expectedError: "",
		},
		{
			name:          "Moderator can moderate user",
			editorRole:    model.RoleModerator,
			targetRole:    model.RoleUser,
			expectedError: "",
		},
		{
			name:          "Moderator cannot moderate admin",
			editorRole:    model.RoleModerator,
			targetRole:    model.RoleAdmin,
			expectedError: "you do not have permission to perform this action",
		},
		{
			name:          "User cannot moderate anyone",
			editorRole:    model.RoleUser,
			targetRole:    model.RoleUser,
			expectedError: "you do not have permission to perform this action",
		},
		{
			name:          "Error getting target role",
			editorRole:    model.RoleAdmin,
			targetRole:    model.RoleUser,
			repoError:     errors.New("database error"),
			expectedError: "failed to fetch target user role",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockRepository{
				GetRoleByIdFunc: func(userId int) (model.Role, error) {
					return tt.targetRole, tt.repoError
				},
			}

			moderService := NewModerService(repo)
			err := moderService.validateModeration(tt.editorRole, 2)

			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError)
			}
		})
	}
}
