package events

type UserUpdated struct {
	UserID    int    `json:"post_id"`
	AvatarURL string `json:"file_url"`
	OldURL    string `json:"old_url"`
}

type UserDeleted struct {
	UserID   int    `json:"post_id"`
	ImageURL string `json:"image_url"`
}
