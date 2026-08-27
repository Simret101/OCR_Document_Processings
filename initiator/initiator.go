package initiator

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"aidoc/docs"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

func Initiate() {
	ctx := context.Background()

	log, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Unable to start logger")
	}

	configName := "config"
	if os.Getenv("CONFIG_NAME") != "" {
		configName = os.Getenv("CONFIG_NAME")
	}

	err = InitConfig(Config{
		Logger: log,
		Names:  []string{configName},
		Path:   "config",
	})
	if err != nil {
		log.Fatal("Unable to start config.", zap.Error(err))
	}

	log.Info("initializing config completed")
	docs.SwaggerInfo.Title = "Document Ocr Processing"
	docs.SwaggerInfo.Description = "API documentation for Ocr Document Processing."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/"
	logger := InitLogger()

	pgpool := initDatabases(logger)

	redisURL := viper.GetString("redis.url")
	redis := initRedis(redisURL, logger)
	minioConfig := MinioConfig{
		Endpoint:  viper.GetString("minio.endpoint"),
		AccessKey: viper.GetString("minio.access_key"),
		SecretKey: viper.GetString("minio.secret_key"),
		UseSSL:    viper.GetBool("minio.ssl"),
		Bucket:    viper.GetString("minio.bucket"),
	}

	minioClient := initMinio(
		minioConfig,
		logger,
	)

	persistance := initPersistance(pgpool, redis, logger, minioClient)

	service := initService(persistance, minioClient, logger)

	transport := initHandler(service, logger)

	server := chi.NewRouter()
	initRoute(server, logger, transport)
	server.Get("/swagger/*", httpSwagger.WrapHandler)

	logger.Info(ctx, "done initializing financial service")

	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", viper.GetString("app.host"), viper.GetInt("app.port")),
		Handler:           server,
		ReadHeaderTimeout: viper.GetDuration("app.timeout"),
		IdleTimeout:       30 * time.Minute,
	}

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT)
		<-sigint
		log.Fatal("HTTP server Shutdown")
	}()

	logger.Info(ctx, fmt.Sprintf("http server listening on port : %d", viper.GetInt("app.port")))

	err = srv.ListenAndServe()
	if err != nil {
		logger.Fatal(ctx, fmt.Sprintf("Could not start HTTP server: %s", err))
	}
}
