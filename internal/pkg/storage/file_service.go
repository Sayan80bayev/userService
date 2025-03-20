package storage

import "mime/multipart"

type FileService interface {
	UploadFile(file multipart.File, header *multipart.FileHeader) (string, error)
	DeleteFileByURL(fileURL string) error
}
