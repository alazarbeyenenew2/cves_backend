package nvd

import (
	"context"
	"encoding/json"
	"time"

	"github.com/alazarbeyenenew2/cves_backend/internal/constant/db/dbinterface"
	"github.com/alazarbeyenenew2/cves_backend/internal/constant/db/generated"
	"github.com/alazarbeyenenew2/cves_backend/internal/constant/dto"
	storage "github.com/alazarbeyenenew2/cves_backend/internal/storege"
	"github.com/alazarbeyenenew2/cves_backend/platform/logger"
	"github.com/jackc/pgtype"
	"github.com/shopspring/decimal"
)

type nvd struct {
	dbinterface *dbinterface.PersistenceDB
	logger      logger.Logger
}

func New(dbinterface *dbinterface.PersistenceDB, logger logger.Logger) storage.CVEStorage {
	return &nvd{
		dbinterface: dbinterface,
		logger:      logger,
	}
}

func (n *nvd) UpdateSyncStatus(ctx context.Context, syncStatus string) {

}
func (n *nvd) GetMeta(ctx context.Context) dto.Meta {
	return dto.Meta{}
}
func (n *nvd) SetSyncStatus(ctx context.Context, syncStatus string) {

}
func (n *nvd) UpsertBatch(ctx context.Context, req []dto.NVDResponse) (int, error) {
	for _, cveBatch := range req {
		for _, cve := range cveBatch.Vulnerabilities {
			metricsScore := 0.0
			baseSeverity := ""
			attackVector := ""
			vectorString := ""
			metrics := cve.CVE.Metrics
			if len(metrics.CvssMetricV31) > 0 {
				metricsScore = metrics.CvssMetricV31[0].CvssData.BaseScore
				baseSeverity = metrics.CvssMetricV31[0].CvssData.BaseSeverity
				attackVector = metrics.CvssMetricV31[0].CvssData.AttackVector
				vectorString = metrics.CvssMetricV31[0].CvssData.VectorString
			} else if len(metrics.CvssMetricV30) > 0 {
				metricsScore = metrics.CvssMetricV30[0].CvssData.BaseScore
				baseSeverity = metrics.CvssMetricV30[0].CvssData.BaseSeverity
				attackVector = metrics.CvssMetricV30[0].CvssData.AttackVector
				vectorString = metrics.CvssMetricV30[0].CvssData.VectorString
			} else if len(metrics.CvssMetricV2) > 0 {
				metricsScore = metrics.CvssMetricV2[0].CvssData.BaseScore
				baseSeverity = metrics.CvssMetricV2[0].BaseSeverity
				attackVector = metrics.CvssMetricV2[0].CvssData.AccessVector
				vectorString = metrics.CvssMetricV2[0].CvssData.VectorString
			}

			year := int16(partTime(cve.CVE.Published).Year())
			metricData, err := json.Marshal(cve.CVE.Metrics)
			if err != nil {
				n.logger.Named("parsing").Warn(ctx, err.Error())
			}
			affected, err := json.Marshal(cve.CVE.Affected)
			if err != nil {
				n.logger.Named("parsing").Warn(ctx, err.Error())
			}

			description, err := json.Marshal(cve.CVE.Descriptions)
			if err != nil {
				n.logger.Named("parsing").Warn(ctx, err.Error())
			}
			data, err := json.Marshal(cve.CVE.References)
			if err != nil {
				n.logger.Named("parsing").Warn(ctx, err.Error())
			}
			weakness, err := json.Marshal(cve.CVE.Weaknesses)
			if err != nil {
				n.logger.Named("parsing").Warn(ctx, err.Error())
			}
			Configurations, err := json.Marshal(cve.CVE.Configurations)
			if err != nil {
				n.logger.Named("parsing").Warn(ctx, err.Error())
			}
			if err := n.dbinterface.Queries.UpsertCVEs(ctx, generated.UpsertCVEsParams{
				CveID:       cve.CVE.ID,
				PublishedAt: partTime(cve.CVE.Published),
				ModifiedAt:  partTime(cve.CVE.LastModified),
				Reference: pgtype.JSONB{Bytes: data, Status: func() pgtype.Status {
					if data != nil {
						return pgtype.Present
					}
					return pgtype.Null
				}()},
				CvssScore:    decimal.NewFromFloat(metricsScore),
				Severity:     &baseSeverity,
				AttackVector: &attackVector,
				VectorString: &vectorString,
				Metrics:      pgtype.JSONB{Bytes: metricData, Status: pgtype.Present},
				Year:         &year,
				Affected:     pgtype.JSONB{Bytes: affected, Status: pgtype.Present},
				CvssVersion:  &cveBatch.Version,
				Descriptions: pgtype.JSONB{Bytes: description, Status: func() pgtype.Status {
					if description != nil {
						return pgtype.Present
					}
					return pgtype.Null
				}()},
				Weaknesses: pgtype.JSONB{Bytes: weakness, Status: func() pgtype.Status {
					if weakness != nil {
						return pgtype.Present
					}
					return pgtype.Null
				}()},
				Configurations: pgtype.JSONB{Bytes: Configurations, Status: func() pgtype.Status {
					if Configurations != nil {
						return pgtype.Present
					}
					return pgtype.Null
				}()},
			},
			); err != nil {
				n.logger.Named("db").Error(ctx, err.Error())
			}
		}
	}

	return len(req), nil
}

func partTime(times string) time.Time {
	layout := "2006-01-02T15:04:05.000"

	t, _ := time.Parse(layout, times)
	return t
}
