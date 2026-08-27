package persistance

import (
	"aidoc/internal/model/db"
	"aidoc/platform/logger"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Persistance struct {
	*db.Queries
	pool *pgxpool.Pool

	log logger.Logger
}

func New(pool *pgxpool.Pool, log logger.Logger) Persistance {
	return Persistance{
		Queries: db.New(pool),
		pool:    pool,
		log:     log,
	}
}
