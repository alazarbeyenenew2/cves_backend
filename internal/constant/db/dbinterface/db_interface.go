package dbinterface

import (
	"github.com/alazarbeyenenew2/cves_backend/internal/constant/db/generated"
	"github.com/alazarbeyenenew2/cves_backend/platform/logger"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PersistenceDB struct {
	*generated.Queries
	Pool *pgxpool.Pool
	log  logger.Logger
}

type Sibling string

func New(pool *pgxpool.Pool, log logger.Logger) PersistenceDB {
	return PersistenceDB{
		Queries: generated.New(pool),
		Pool:    pool,
		log:     log,
	}
}
