package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/abijith/moneymate-backend/services/support/config"
	"github.com/abijith/moneymate-backend/services/support/internal/app"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to Load Config: %v", err)
	}

	supportApp, err := app.Build(cfg)
	if err != nil {
		log.Fatalf("Failed to build app: %v", err)
	}
	defer supportApp.Close()

	go func() {
		if err := supportApp.Run(); err != nil {
			log.Fatalf("App run failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("Shutdown signal received, gracefully shutting down...")
	supportApp.Close()
	log.Println("Support service stopped cleanly ")
}
