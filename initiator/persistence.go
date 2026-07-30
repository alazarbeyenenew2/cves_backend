package initiator

import (
	"github.com/alazarbeyenenew2/cves_backend/internal/constant/db/dbinterface"
	"github.com/alazarbeyenenew2/cves_backend/platform/logger"
	"github.com/go-redis/redis/v8"
)

type Persistance struct {
}

func initPersistence(persistencedb *dbinterface.PersistenceDB, redisClient *redis.Client, log logger.Logger) *Persistance {
	return &Persistance{}
}
