package service

import (
	"context"
	"encoding/json"
	"fmt"
	"userService/internal/model"

	"github.com/Sayan80bayev/go-project/pkg/logging"
	storage "github.com/Sayan80bayev/go-project/pkg/objectStorage"
	"userService/internal/events"
)

var logger = logging.GetLogger()

func CreateUserHandler(repository UserRepository) func(data json.RawMessage) error {
	return func(data json.RawMessage) error {
		ctx := context.WithoutCancel(context.Background())

		var e events.UserCreatedPayload
		if err := json.Unmarshal(data, &e); err != nil {
			return fmt.Errorf("failed to unmarshal UserCreatedPayload: %w", err)
		}

		needsCompletion := false
		if e.Firstname == "" || e.Lastname == "" || e.Firstname == "null" || e.Lastname == "null" {
			needsCompletion = true
		}

		user := &model.User{
			ID:              e.UserID,
			Firstname:       e.Firstname,
			Lastname:        e.Lastname,
			Email:           e.Email,
			NeedsCompletion: needsCompletion,
		}

		if err := repository.CreateUser(ctx, user); err != nil {
			logger.Error("failed to create user", err)
		}

		logger.Infof("Created user profile: %+v", user)
		return nil
	}
}

// UserUpdatedHandler handles user update events
func UserUpdatedHandler(fileStorage storage.FileStorage) func(data json.RawMessage) error {
	return func(data json.RawMessage) error {
		ctx := context.WithoutCancel(context.Background())
		var e events.UserUpdatedPayload
		if err := json.Unmarshal(data, &e); err != nil {
			return fmt.Errorf("failed to unmarshal UserUpdatedPayload: %w", err)
		}

		if e.OldURL != "" && e.OldURL != e.AvatarURL {
			if err := fileStorage.DeleteFileByURL(ctx, e.OldURL); err != nil {
				logger.Errorf("Error deleting old file on user update: %v", err)
			}
		}
		return nil
	}
}

// UserDeletedHandler handles user deletion events
func UserDeletedHandler(fileStorage storage.FileStorage) func(data json.RawMessage) error {
	return func(data json.RawMessage) error {
		ctx := context.WithoutCancel(context.Background())

		var e events.UserDeletedPayload
		if err := json.Unmarshal(data, &e); err != nil {
			return fmt.Errorf("failed to unmarshal UserDeletedPayload: %w", err)
		}

		if err := fileStorage.DeleteFileByURL(ctx, e.ImageURL); err != nil {
			logger.Errorf("Error deleting file on user delete: %v", err)
		}
		return nil
	}
}
