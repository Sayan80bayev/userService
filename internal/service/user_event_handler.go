package service

import (
	"encoding/json"
	"fmt"

	"github.com/Sayan80bayev/go-project/pkg/logging"
	storage "github.com/Sayan80bayev/go-project/pkg/objectStorage"
	"userService/internal/events"
)

var logger = logging.GetLogger()

// UserUpdatedHandler handles user update events
func UserUpdatedHandler(fileStorage storage.FileStorage) func(data json.RawMessage) error {
	return func(data json.RawMessage) error {
		var e events.UserUpdated
		if err := json.Unmarshal(data, &e); err != nil {
			return fmt.Errorf("failed to unmarshal UserUpdated: %w", err)
		}

		if e.OldURL != "" && e.OldURL != e.AvatarURL {
			if err := fileStorage.DeleteFileByURL(e.OldURL); err != nil {
				logger.Errorf("Error deleting old file on user update: %v", err)
			}
		}
		return nil
	}
}

// UserDeletedHandler handles user deletion events
func UserDeletedHandler(fileStorage storage.FileStorage) func(data json.RawMessage) error {
	return func(data json.RawMessage) error {
		var e events.UserDeleted
		if err := json.Unmarshal(data, &e); err != nil {
			return fmt.Errorf("failed to unmarshal UserDeleted: %w", err)
		}

		if err := fileStorage.DeleteFileByURL(e.ImageURL); err != nil {
			logger.Errorf("Error deleting file on user delete: %v", err)
		}
		return nil
	}
}
