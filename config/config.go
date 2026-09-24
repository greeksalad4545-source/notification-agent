package config

import "os"

type Config struct {
	ServiceBusNamespace        string
	ServiceBusQueue            string
	ServiceBusConnectionString string
	DatabaseURL                string
}

func Load() Config {
	return Config{
		ServiceBusNamespace:        os.Getenv("SERVICE_BUS_NAMESPACE"),
		ServiceBusQueue:            os.Getenv("SERVICE_BUS_QUEUE"),
		ServiceBusConnectionString: os.Getenv("SERVICE_BUS_CONNECTION_STRING"),
		DatabaseURL:                os.Getenv("DATABASE_URL"),
	}
}
