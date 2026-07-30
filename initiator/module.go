package initiator

import (
	"github.com/alazarbeyenenew2/cves_backend/platform/logger"
	"github.com/alazarbeyenenew2/cves_backend/platform/workerpool"
)

type Module struct {
}

func initModule(
	persistence *Persistance,
	log logger.Logger,
	pool *workerpool.WorkerPool,
) *Module {
	return &Module{}
}
