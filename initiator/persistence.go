package initiator

import (
	"github.com/alazarbeyenenew2/cves_backend/internal/constant/db/dbinterface"
	storage "github.com/alazarbeyenenew2/cves_backend/internal/storege"
	"github.com/alazarbeyenenew2/cves_backend/internal/storege/nvd"
	"github.com/alazarbeyenenew2/cves_backend/platform/logger"
	"github.com/go-redis/redis/v8"
)

type Persistance struct {
	NVDStorage storage.CVEStorage
}

func initPersistence(persistencedb *dbinterface.PersistenceDB, redisClient *redis.Client, log logger.Logger) *Persistance {
	return &Persistance{
		NVDStorage: nvd.New(persistencedb, log),
	}
}
