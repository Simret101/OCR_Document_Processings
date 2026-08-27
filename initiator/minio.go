package initiator

import (
	"context"
	"time"

	"aidoc/platform/logger"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
	Bucket    string
}

func initMinio(cfg MinioConfig, log logger.Logger) *minio.Client {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(
			cfg.AccessKey,
			cfg.SecretKey,
			"",
		),
		Secure: cfg.UseSSL,
	})

	if err != nil {
		log.Fatal(context.Background(), "failed to initialize minio client", zap.Error(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.ListBuckets(ctx)
	if err != nil {
		log.Fatal(context.Background(), "failed to connect to minio", zap.Error(err))
	}

	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		log.Fatal(context.Background(), "failed checking bucket", zap.Error(err))
	}

	if !exists {
		err = client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatal(context.Background(), "failed creating bucket", zap.Error(err))
		}
	}

	log.Info(context.Background(), "minio connected successfully",
		zap.String("endpoint", cfg.Endpoint),
		zap.String("bucket", cfg.Bucket),
	)

	return client
}
