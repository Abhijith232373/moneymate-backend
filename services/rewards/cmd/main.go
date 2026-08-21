package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/moneymate-2026/moneymate-backend/services/rewards/config"
	"github.com/moneymate-2026/moneymate-backend/services/rewards/internal/app"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	rewardsApp, err := app.Build(cfg)
	if err != nil {
		log.Fatalf("failed to build rewards app: %v", err)
	}
	defer rewardsApp.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if rewardsApp.KafkaConsumer != nil {
		go rewardsApp.KafkaConsumer.Run(ctx, rewardsApp.HandleKafkaMessage)
	}

	go func() {
		if err := rewardsApp.Run(); err != nil {
			log.Fatalf("rewards app failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("shutdown signal received, stopping rewards service")
	cancel()
	rewardsApp.Close()
	log.Println("rewards service stopped cleanly")
}
