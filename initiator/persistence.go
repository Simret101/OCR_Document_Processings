package initiator

import (
	"aidoc/internal/repository"
	"aidoc/internal/repository/imagestore"
	"aidoc/platform/logger"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/spf13/viper"
)

type Persistence struct {
	ImageStore repository.Image
}

func initPersistance(pool *pgxpool.Pool, redis *redis.Client, log logger.Logger, minioClient *minio.Client) *Persistence {
	image := imagestore.NewMinioStorage(minioClient,viper.GetString("bucket"))

	return &Persistence{
		ImageStore: image,
	}
}
