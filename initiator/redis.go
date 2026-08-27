package initiator

import (
	"context"
	"time"

	"aidoc/platform/logger"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func initRedis(url string, log logger.Logger) *redis.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts, err := redis.ParseURL(url)
	if err != nil {
		log.Fatal(ctx, "Invalid redis url", zap.Error(err))
	}

	client := redis.NewClient(opts)
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatal(ctx, "Unable to connect to redis", zap.Error(err))
	}

	log.Info(ctx, "Redis initialized")
	return client
}
