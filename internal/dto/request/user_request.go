package request

import "mime/multipart"

type UserRequest struct {
	Username    string `json:"username"`
	About       string `json:"about"`
	DateOfBirth string `json:"dateOfBirth"`
	Avatar      multipart.File
	Header      *multipart.FileHeader
}
