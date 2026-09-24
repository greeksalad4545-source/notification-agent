package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"notification-agent/agent"
	"notification-agent/config"
	"notification-agent/servicebus"
	"notification-agent/storage"
)

func main() {
	cfg := config.Load()

	fmt.Println("Configuration loaded")
	fmt.Println("Service Bus Queue:", cfg.ServiceBusQueue)

	// Create a context that is cancelled when the application
	// receives Ctrl+C or a termination signal.
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// Connect to PostgreSQL
	db, err := storage.NewPostgres(ctx, cfg.DatabaseURL)
	if err != nil {
		fmt.Println("Failed to connect to PostgreSQL:", err)
		return
	}
	defer db.Close(context.Background())

	err = db.Ping(ctx)
	if err != nil {
		fmt.Println("PostgreSQL ping failed:", err)
		return
	}

	fmt.Println("Connected to PostgreSQL successfully!")

	// Create database tables
	err = db.CreateTables(ctx)
	if err != nil {
		fmt.Println("Failed to create database tables:", err)
		return
	}

	fmt.Println("Database tables ready.")

	// Connect to Azure Service Bus
	consumer, err := servicebus.NewConsumer(
		cfg.ServiceBusConnectionString,
		cfg.ServiceBusQueue,
	)
	if err != nil {
		fmt.Println("Failed to connect to Service Bus:", err)
		return
	}
	defer consumer.Close(context.Background())

	fmt.Println("Connected to Azure Service Bus")

	// Create notification agent
	notificationAgent := agent.New(
		consumer,
		db,
	)

	fmt.Println("Notification Agent starting...")

	// Start the agent
	err = notificationAgent.Run(ctx)
	if err != nil {
		fmt.Println("Notification Agent stopped with error:", err)
		return
	}

	fmt.Println("Notification Agent stopped.")
}
