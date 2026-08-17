package initiator

import (
	"github.com/alazarbeyenenew2/cves_backend/internal/module"
	"github.com/alazarbeyenenew2/cves_backend/internal/module/nvd"
	"github.com/alazarbeyenenew2/cves_backend/platform/logger"
	"github.com/alazarbeyenenew2/cves_backend/platform/workerpool"
	"github.com/spf13/viper"
)

type Module struct {
	NVD      module.NVD
	platform Platform
}

func initModule(
	persistence *Persistance,
	log logger.Logger,
	pool *workerpool.WorkerPool,
	platform *Platform,
) *Module {
	return &Module{
		NVD: nvd.New(persistence.NVDStorage, &platform.NVDClient, viper.GetDuration("nvd.interval"), viper.GetInt("nvd.start_year"), log),
	}
}
