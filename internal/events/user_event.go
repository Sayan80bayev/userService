package events

const (
	UserUpdated = "UserUpdated"
	UserDeleted = "UserDeleted"
)

type UserUpdatedPayload struct {
	UserID    int    `json:"post_id"`
	AvatarURL string `json:"file_url"`
	OldURL    string `json:"old_url"`
}

type UserDeletedPayload struct {
	UserID   int    `json:"post_id"`
	ImageURL string `json:"image_url"`
}
