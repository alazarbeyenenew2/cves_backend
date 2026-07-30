package initiator

import (
	"fmt"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Names  []string
	Path   string
	Type   string
	Logger *zap.Logger
}

func InitConfig(config Config) error {
	if config.Logger == nil {
		return fmt.Errorf("logger cannot be nil")
	}

	if len(config.Names) == 0 || config.Path == "" {
		return fmt.Errorf("config names and path cannot be empty")
	}

	if config.Type == "" {
		config.Type = "yaml"
	}

	viper.AddConfigPath(config.Path)
	viper.SetConfigType(config.Type)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("APPLICATION")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	for i, name := range config.Names {
		viper.SetConfigName(name)

		var err error
		if i == 0 {
			err = viper.ReadInConfig()
		} else {
			err = viper.MergeInConfig()
		}

		if err != nil {
			return fmt.Errorf("failed to load config %s: %w", name, err)
		}
	}

	if err := validateConfig(); err != nil {
		return fmt.Errorf("configuration validation failed: %v", err)
	}

	if err := setupConfigWatcher(config.Logger); err != nil {
		config.Logger.Warn("Failed to setup config watcher",
			zap.Error(err),
		)
	}

	config.Logger.Info("Configuration initialized successfully",
		zap.Strings("config_files", config.Names),
		zap.String("config_path", config.Path),
	)

	return nil
}

func validateConfig() error {

	requiredKeys := []string{}

	for _, key := range requiredKeys {
		if !viper.IsSet(key) || viper.GetString(key) == "" {
			return fmt.Errorf("missing required configuration key: %s", key)
		}
	}
	return nil
}

func setupConfigWatcher(log *zap.Logger) error {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Info("Configuration file changed",
			zap.String("file", e.Name),
			zap.String("operation", e.Op.String()),
		)

		if err := validateConfig(); err != nil {
			log.Error("Invalid configuration after change",
				zap.Error(err),
			)
		}
	})
	return nil
}
