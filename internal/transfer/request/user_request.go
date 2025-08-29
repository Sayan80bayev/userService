package request

import "mime/multipart"

// UserRequest is the incoming DTO for creating/updating a user profile
type UserRequest struct {
	Email       string   `form:"email" binding:"required,email"`
	Lastname    string   `form:"lastname" binding:"required,min=3,max=20"`
	Firstname   string   `form:"firstname" binding:"required,min=3,max=20"`
	About       string   `form:"about,omitempty" validate:"omitempty,max=500"`
	DateOfBirth string   `form:"dateOfBirth,omitempty"` // keep as string, parse to time.Time later
	Gender      string   `form:"gender,omitempty" validate:"omitempty,oneof=male female other"`
	Location    string   `form:"location,omitempty" validate:"omitempty,max=100"`
	Socials     []string `form:"socials[]" validate:"omitempty,dive,url"`

	// File upload fields (unchanged)
	Avatar multipart.File
	Header *multipart.FileHeader
}
