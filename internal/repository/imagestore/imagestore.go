package imagestore

import (
	"context"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
)

type MinioStorage struct {
	client *minio.Client
	bucket string
}

func NewMinioStorage(client *minio.Client, bucket string) *MinioStorage {
	return &MinioStorage{
		client: client,
		bucket: bucket,
	}
}

func (m *MinioStorage) Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {

	_, err := m.client.PutObject(ctx, m.bucket, objectName, reader, size, minio.PutObjectOptions{ContentType: contentType})

	if err != nil {
		return "", err
	}

	return objectName, nil
}

func (m *MinioStorage) Download(ctx context.Context, objectName string) (io.ReadCloser, error) {

	return m.client.GetObject(ctx, m.bucket, objectName, minio.GetObjectOptions{})
}

func (m *MinioStorage) Delete(ctx context.Context, objectName string) error {

	return m.client.RemoveObject(ctx, m.bucket, objectName, minio.RemoveObjectOptions{})
}
func (m *MinioStorage) GetPresignedURL(ctx context.Context, objectName string) (string, error) {

	url, err := m.client.PresignedGetObject(ctx, m.bucket, objectName, 15*time.Minute, nil)

	if err != nil {
		return "", err
	}

	return url.String(), nil
}
