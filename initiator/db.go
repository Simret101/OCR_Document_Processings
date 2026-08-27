package initiator

import (
	"context"
	"fmt"
	"time"

	"aidoc/platform/logger"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
)

func initDatabases(log logger.Logger) *pgxpool.Pool {
	pgPool := initPostgres("db.url", "_", log)

	return pgPool
}

func initPostgres(configKey, dbName string, log logger.Logger) *pgxpool.Pool {
	url := viper.GetString(configKey)
	if url == "" {
		log.Fatal(context.Background(), "database url is empty")
	}

	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatal(context.Background(), "unable to parse postgres config")
	}

	config.MaxConnIdleTime = viper.GetDuration("database.idle_conn_timeout")
	if config.MaxConnIdleTime == 0 {
		config.MaxConnIdleTime = 4 * time.Minute
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal(context.Background(), fmt.Sprintf("failed to connect to postgres (%s)", dbName))
	}

	log.Info(context.Background(), "PostgreSQL connected successfully")
	return pool
}
