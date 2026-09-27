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
	KeyVaultName               string
	SendGridFrom               string
}

func Load() (Config, error) {
	cfg := Config{
		ServiceBusNamespace:        os.Getenv("SERVICE_BUS_NAMESPACE"),
		ServiceBusQueue:            os.Getenv("SERVICE_BUS_QUEUE"),
		ServiceBusConnectionString: os.Getenv("SERVICE_BUS_CONNECTION_STRING"),
		DatabaseURL:                os.Getenv("DATABASE_URL"),
		KeyVaultName:               os.Getenv("KEY_VAULT_NAME"),
		SendGridFrom:               os.Getenv("SENDGRID_FROM"),
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

	if cfg.KeyVaultName == "" {
		return Config{}, fmt.Errorf("KEY_VAULT_NAME is required")
	}

	if cfg.SendGridFrom == "" {
		return Config{}, fmt.Errorf("SENDGRID_FROM is required")
	}

	return cfg, nil
}
