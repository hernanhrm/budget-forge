package localconfig

import (
	"fmt"
	"os"

	shared_domain "github.com/hernanhrm/budget-forge/pkg/shared_domain"
)

func GetConfig(log shared_domain.Logger) (LocalConfig, error) {
	config := LocalConfig{
		Service: Service{
			Port:  getPort,
			Name:  getEnvAsString("SERVICE_NAME"),
			Debug: getDebug(),
		},
		Database: Database{
			URL: getEnvAsString("DATABASE_URL"),
		},
	}

	if config.Database.URL == "" {
		return LocalConfig{}, fmt.Errorf("environment variable DATABASE_URL is required")
	}

	log.Debug("configuration loaded",
		"service_name", config.Service.Name,
		"service_port", config.Service.Port(),
	)

	return config, nil
}

func getEnvAsString(key string) string {
	return os.Getenv(key)
}
