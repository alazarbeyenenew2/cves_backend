package initiator

import (
	"github.com/alazarbeyenenew2/cves_backend/internal/constant/dto"
	"github.com/alazarbeyenenew2/cves_backend/platform/logger"
	NVDClient "github.com/alazarbeyenenew2/cves_backend/platform/nvd_client"
	"github.com/spf13/viper"
)

type Platform struct {
	NVDClient NVDClient.NVDClient
}

func initPlatform(logger logger.Logger) *Platform {
	return &Platform{
		NVDClient: *NVDClient.NewNVDClient(dto.NVDClient{
			ApiKey:     viper.GetString("nvd.api_key"),
			BaseURL:    viper.GetString("nvd.base_url"),
			PageSize:   viper.GetInt("nvd.page_size"),
			MaxRetries: viper.GetInt("nvd.max_retries"),
			Logger:     logger,
		}),
	}
}
