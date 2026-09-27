package config

import (
	"os"
	"testing"
)

func clearConfigEnv(t *testing.T) {
	t.Helper()

	os.Unsetenv("SERVICE_BUS_NAMESPACE")
	os.Unsetenv("SERVICE_BUS_QUEUE")
	os.Unsetenv("SERVICE_BUS_CONNECTION_STRING")
	os.Unsetenv("DATABASE_URL")
}

func TestLoadSuccess(t *testing.T) {
	clearConfigEnv(t)
	defer clearConfigEnv(t)

	os.Setenv("SERVICE_BUS_NAMESPACE", "test-namespace")
	os.Setenv("SERVICE_BUS_QUEUE", "test-queue")
	os.Setenv("SERVICE_BUS_CONNECTION_STRING", "test-connection-string")
	os.Setenv("DATABASE_URL", "postgres://test")

	cfg, err := Load()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.ServiceBusQueue != "test-queue" {
		t.Errorf("expected queue test-queue, got %s", cfg.ServiceBusQueue)
	}

	if cfg.DatabaseURL != "postgres://test" {
		t.Errorf("expected database URL postgres://test, got %s", cfg.DatabaseURL)
	}
}

func TestLoadMissingQueue(t *testing.T) {
	clearConfigEnv(t)
	defer clearConfigEnv(t)

	os.Setenv("SERVICE_BUS_CONNECTION_STRING", "test-connection-string")
	os.Setenv("DATABASE_URL", "postgres://test")

	_, err := Load()

	if err == nil {
		t.Fatal("expected error when SERVICE_BUS_QUEUE is missing")
	}
}

func TestLoadMissingServiceBusConnectionString(t *testing.T) {
	clearConfigEnv(t)
	defer clearConfigEnv(t)

	os.Setenv("SERVICE_BUS_QUEUE", "test-queue")
	os.Setenv("DATABASE_URL", "postgres://test")

	_, err := Load()

	if err == nil {
		t.Fatal("expected error when SERVICE_BUS_CONNECTION_STRING is missing")
	}
}

func TestLoadMissingDatabaseURL(t *testing.T) {
	clearConfigEnv(t)
	defer clearConfigEnv(t)

	os.Setenv("SERVICE_BUS_QUEUE", "test-queue")
	os.Setenv("SERVICE_BUS_CONNECTION_STRING", "test-connection-string")

	_, err := Load()

	if err == nil {
		t.Fatal("expected error when DATABASE_URL is missing")
	}
}
