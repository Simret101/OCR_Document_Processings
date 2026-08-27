package preprocessing

import (
	"aidoc/platform/ocr/preprocessing/compression"
	"context"
	"image"
	"io"

	"github.com/minio/minio-go/v7"
)

type Storage struct {
	client       *minio.Client
	bucket       string
	width        int
	height       int
	maxDimension int
	clipLimit    float64
	tileGridSize image.Point
	KernelSize   int
	compressor   compression.Compressor
}

func NewStorage(client *minio.Client, bucket string, width int, height int, maxDimension int, clipLimit float64, tileGridSize image.Point, KernelSize int,compressor   compression.Compressor) *Storage {
	return &Storage{
		client:       client,
		bucket:       bucket,
		width:        width,
		height:       height,
		maxDimension: maxDimension,
		clipLimit:    clipLimit,
		tileGridSize: tileGridSize,
		KernelSize:   KernelSize,
		compressor: compressor,
	}
}
func (s *Storage) GetImage(ctx context.Context, objectName string) ([]byte, error) {
	object, err := s.client.GetObject(ctx, s.bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer object.Close()

	return io.ReadAll(object)
}
