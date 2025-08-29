package events

import "github.com/google/uuid"

const (
	UserCreated = "UserCreated"
	UserUpdated = "UserUpdated"
	UserDeleted = "UserDeleted"
)

type UserCreatedPayload struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
}
type UserUpdatedPayload struct {
	UserID    uuid.UUID `json:"user_id"`
	AvatarURL string    `json:"file_url"`
	OldURL    string    `json:"old_url"`
}

type UserDeletedPayload struct {
	UserID   uuid.UUID `json:"user_id"`
	ImageURL string    `json:"image_url"`
}
