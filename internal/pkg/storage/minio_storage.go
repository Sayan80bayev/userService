package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"mime/multipart"
	"strings"
	"userService/internal/config"
	"userService/pkg/logging"
)

var logger = logging.GetLogger()

func Init(cfg *config.Config) *minio.Client {
	minioClient, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		logger.Error("Failed to initialize MinIO client:", err)
		return nil
	}
	logger.Info("Successfully connected to MinIO")
	return minioClient
}

func UploadFile(file multipart.File, header *multipart.FileHeader, cfg *config.Config, minioClient *minio.Client) (string, error) {
	if file == nil || header == nil {
		return "", errors.New("invalid file")
	}
	defer file.Close()

	objectName := header.Filename
	contentType := header.Header.Get("Content-Type")
	bucketName := cfg.MinioBucket

	_, err := minioClient.PutObject(
		context.Background(),
		bucketName,
		objectName,
		file,
		header.Size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	fileURL := fmt.Sprintf("http://%s/%s/%s", cfg.MinioEndpoint, bucketName, objectName)

	return fileURL, nil
}

func DeleteFileByURL(fileURL string, cfg *config.Config, minioClient *minio.Client) error {
	if fileURL == "" {
		return errors.New("missing file_url parameter")
	}

	prefix := fmt.Sprintf("http://%s/", cfg.MinioEndpoint)
	if !strings.HasPrefix(fileURL, prefix) {
		return errors.New("invalid file_url format")
	}
	path := strings.TrimPrefix(fileURL, prefix)

	parts := strings.SplitN(path, "/", 2)
	if len(parts) < 2 {
		return errors.New("invalid file_url format")
	}

	bucketName, objectName := parts[0], parts[1]

	err := minioClient.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

type MinioStorage struct {
	client *minio.Client
	cfg    *config.Config
}

func NewMinioStorage(client *minio.Client, cfg *config.Config) *MinioStorage {
	return &MinioStorage{client: client, cfg: cfg}
}

func (s *MinioStorage) UploadFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	return UploadFile(file, header, s.cfg, s.client)
}

func (s *MinioStorage) DeleteFileByURL(fileURL string) error {
	return DeleteFileByURL(fileURL, s.cfg, s.client)
}
