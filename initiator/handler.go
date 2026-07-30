package initiator

import (
	"github.com/alazarbeyenenew2/cves_backend/platform/logger"
	"github.com/alazarbeyenenew2/cves_backend/platform/workerpool"
)

type Handler struct{}

func initHandler(module *Module, log logger.Logger, wp *workerpool.WorkerPool) *Handler {
	return &Handler{}
}
