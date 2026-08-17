package nvd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/alazarbeyenenew2/cves_backend/internal/module"
	storage "github.com/alazarbeyenenew2/cves_backend/internal/storege"
	"github.com/alazarbeyenenew2/cves_backend/platform/logger"
	"github.com/alazarbeyenenew2/cves_backend/platform/nvd_client"
)

type Scheduler struct {
	CVEStore  storage.CVEStorage
	client    *nvd_client.NVDClient
	logger    logger.Logger
	stopCh    chan struct{}
	startYear int
	interval  time.Duration
}

func New(s storage.CVEStorage, client *nvd_client.NVDClient, interval time.Duration, startYear int, logger logger.Logger) module.NVD {
	return &Scheduler{
		CVEStore:  s,
		client:    client,
		stopCh:    make(chan struct{}),
		startYear: startYear,
		interval:  interval,
		logger:    logger,
	}
}

func (sc *Scheduler) Start() {
	go sc.run()
}

func (sc *Scheduler) Stop() {
	close(sc.stopCh)
}

func (sc *Scheduler) run() {
	ticker := time.NewTicker(sc.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			sc.logger.Named("scheduler").Info(context.Background(), "[CRON] Starting scheduled CVE sync...")
			sc.doSync()
		case <-sc.stopCh:
			sc.logger.Named("scheduler").Info(context.Background(), "[CRON] Scheduler stopped.")
			return
		}
	}
}

func (sc *Scheduler) TriggerNow() {
	go sc.doSync()
}

func (sc *Scheduler) doSync() {
	sc.CVEStore.UpdateSyncStatus(context.Background(), "syncing")

	meta := sc.CVEStore.GetMeta(context.Background())

	if meta.LastFetchTime.IsZero() {
		sc.doHistoricalSync()
	} else {
		sc.doIncrementalSync(meta.LastFetchTime.Add(-5 * time.Minute))
	}
}

func (sc *Scheduler) doHistoricalSync() {
	startYear := sc.startYear
	if startYear < 1999 {
		startYear = 1999
	}
	currentYear := time.Now().Year()
	totalNew := 0

	sc.logger.Named("scheduler").Info(context.Background(), fmt.Sprintf("[CRON] Historical sync: fetching %d–%d year by year", startYear, currentYear))

	for year := startYear; year <= currentYear; year++ {
		yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		yearEnd := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)
		if yearEnd.After(time.Now()) {
			yearEnd = time.Now()
		}

		sc.logger.Named("scheduler").Info(context.Background(), fmt.Sprintf("[CRON] Fetching year %d...", year))

		sc.CVEStore.SetSyncStatus(context.Background(), fmt.Sprintf("syncing %d", year))

		cves, err := sc.client.FetchRange(yearStart, yearEnd, func(fetched, total int) {
			sc.logger.Named("scheduler").Info(context.Background(), fmt.Sprintf("[CRON] %d: %d / %d CVEs", year, fetched, total))
		})
		if err != nil {
			sc.logger.Named("scheduler").Info(context.Background(), fmt.Sprintf("[CRON] Error fetching year %d: %v — continuing to next year", year, err))
			continue
		}

		newCount, err := sc.CVEStore.UpsertBatch(context.Background(), cves)
		if err != nil {
			log.Printf("[CRON] Store error for year %d: %v", year, err)
			continue
		}
		totalNew += newCount
		sc.logger.Named("scheduler").Info(context.Background(), fmt.Sprintf("[CRON] Year %d complete: %d fetched, %d new", year, len(cves), newCount))
	}

	log.Printf("[CRON] Historical sync complete. %d total new CVEs stored.", totalNew)

	sc.CVEStore.SetSyncStatus(context.Background(), "idle")
}

func (sc *Scheduler) doIncrementalSync(since time.Time) {
	sc.logger.Named("scheduler").Info(context.Background(), fmt.Sprintf("[CRON] Incremental sync since %s", since.Format(time.RFC3339)))
	cves, err := sc.client.FetchSince(since, func(fetched, total int) {
		log.Printf("[CRON] Progress: %d / %d CVEs fetched", fetched, total)
	})
	if err != nil {
		sc.logger.Named("scheduler").Info(context.Background(), fmt.Sprintf("[CRON] Fetch error: %v", err))
		sc.CVEStore.SetSyncStatus(context.Background(), fmt.Sprintf("error: %v", err))
		return
	}

	newCount, err := sc.CVEStore.UpsertBatch(context.Background(), cves)
	if err != nil {
		sc.logger.Named("scheduler").Info(context.Background(), fmt.Sprintf("[CRON] Store error: %v", err))
		sc.CVEStore.SetSyncStatus(context.Background(), fmt.Sprintf("error: %v", err))
		return
	}

	sc.logger.Named("scheduler").Info(context.Background(), fmt.Sprintf("[CRON] Incremental sync complete. %d fetched, %d new.", len(cves), newCount))

	sc.CVEStore.SetSyncStatus(context.Background(), "idle")
}
