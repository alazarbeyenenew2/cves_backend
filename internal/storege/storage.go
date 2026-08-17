package storage

import (
	"context"

	"github.com/alazarbeyenenew2/cves_backend/internal/constant/dto"
)

type CVEStorage interface {
	UpdateSyncStatus(ctx context.Context, syncStatus string)
	GetMeta(ctx context.Context) dto.Meta
	SetSyncStatus(ctx context.Context, syncStatus string)
	UpsertBatch(ctx context.Context, req []dto.NVDResponse) (int, error)
}
