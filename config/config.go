package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServiceBusNamespace        string
	ServiceBusQueue            string
	ServiceBusConnectionString string
	DatabaseURL                string
}

func Load() (Config, error) {
	cfg := Config{
		ServiceBusNamespace:        os.Getenv("SERVICE_BUS_NAMESPACE"),
		ServiceBusQueue:            os.Getenv("SERVICE_BUS_QUEUE"),
		ServiceBusConnectionString: os.Getenv("SERVICE_BUS_CONNECTION_STRING"),
		DatabaseURL:                os.Getenv("DATABASE_URL"),
	}

	if cfg.ServiceBusQueue == "" {
		return Config{}, fmt.Errorf("SERVICE_BUS_QUEUE is required")
	}

	if cfg.ServiceBusConnectionString == "" {
		return Config{}, fmt.Errorf("SERVICE_BUS_CONNECTION_STRING is required")
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}
