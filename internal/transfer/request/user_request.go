package request

import "mime/multipart"

type UserRequest struct {
	Username    string `form:"username" binding:"required"`
	About       string `form:"about"`
	DateOfBirth string `form:"dateOfBirth"`
	Avatar      multipart.File
	Header      *multipart.FileHeader
}
