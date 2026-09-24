package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"notification-agent/agent"
	"notification-agent/config"
	"notification-agent/grpcserver"
	"notification-agent/proto"
	"notification-agent/servicebus"
	"notification-agent/storage"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	fmt.Println("Configuration loaded")
	fmt.Println("Service Bus Queue:", cfg.ServiceBusQueue)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

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

	err = db.CreateTables(ctx)
	if err != nil {
		fmt.Println("Failed to create database tables:", err)
		return
	}

	fmt.Println("Database tables ready.")

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

	grpcListener, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Println("Failed to start gRPC listener:", err)
		return
	}
	defer grpcListener.Close()

	grpcServer := grpc.NewServer()

	notificationGRPCServer := grpcserver.NewServer(db)

	proto.RegisterNotificationServiceServer(
		grpcServer,
		notificationGRPCServer,
	)

	go func() {
		fmt.Println("gRPC server listening on :50051")

		if err := grpcServer.Serve(grpcListener); err != nil {
			fmt.Println("gRPC server stopped:", err)
		}
	}()

	notificationAgent := agent.New(
		consumer,
		db,
	)

	fmt.Println("Notification Agent starting...")

	err = notificationAgent.Run(ctx)
	if err != nil {
		fmt.Println("Notification Agent stopped with error:", err)
	}

	grpcServer.GracefulStop()

	fmt.Println("Notification Agent stopped.")
}
