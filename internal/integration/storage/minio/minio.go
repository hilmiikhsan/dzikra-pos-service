package minio

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"
)

var _ MinioService = &minioService{}

type minioService struct {
	client     *minio.Client
	bucketName string
}

func NewMinioService(client *minio.Client, bucketName string) *minioService {
	return &minioService{
		client:     client,
		bucketName: bucketName,
	}
}

func (s *minioService) UploadFile(ctx context.Context, objectName string, file multipart.File, fileSize int64, contentType string) (string, error) {
	uploadInfo, err := s.client.PutObject(ctx, s.bucketName, objectName, file, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		log.Error().Err(err).Msg("minio::UploadFile - error uploading file")
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	log.Info().Msgf("Successfully uploaded %s of size %d", uploadInfo.Key, uploadInfo.Size)

	fileURL := fmt.Sprintf("%s/%s/%s", s.client.EndpointURL().String(), s.bucketName, objectName)

	return fileURL, nil
}

func (s *minioService) DeleteFile(ctx context.Context, objectName string) error {
	err := s.client.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		log.Error().Err(err).Msgf("minio::DeleteFile - error deleting object %s", objectName)
		return fmt.Errorf("failed to delete object %s: %w", objectName, err)
	}

	log.Info().Msgf("minio::DeleteFile - successfully deleted object %s", objectName)

	return nil
}
