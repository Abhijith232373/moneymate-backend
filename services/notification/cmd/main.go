package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/moneymate-2026/moneymate-backend/services/notification/config"
	"github.com/moneymate-2026/moneymate-backend/services/notification/internal/app"
	kafkaconsumer "github.com/moneymate-2026/moneymate-backend/services/notification/internal/infra/kafkaconsumer"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	notificationApp, err := app.Build(cfg)
	if err != nil {
		log.Fatalf("Failed to build app: %v", err)
	}
	defer notificationApp.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go notificationApp.KafkaConsumer.Run(ctx, func(ctx context.Context, payload []byte) error {
		return kafkaconsumer.HandleEvent(ctx, notificationApp.NotificationUC, payload)
	})

	go func() {
		if err := notificationApp.Run(); err != nil {
			log.Fatalf("App run failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received, gracefully shutting down...")
	cancel()
	notificationApp.Close()
	log.Println("Notification service stopped cleanly")
}