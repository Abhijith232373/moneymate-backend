package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moneymate-2026/moneymate-backend/services/merchant/config"
	"github.com/moneymate-2026/moneymate-backend/services/merchant/internal/adapter/postgres"
	"github.com/moneymate-2026/moneymate-backend/services/merchant/internal/adapter/postgres/repo"
	transporthttp "github.com/moneymate-2026/moneymate-backend/services/merchant/internal/transport/http"
	ws "github.com/moneymate-2026/moneymate-backend/services/merchant/internal/transport/websocket"
	"github.com/moneymate-2026/moneymate-backend/services/merchant/internal/usecases"
)

type App struct {
	HTTPServer *fiber.App
	DB         *pgxpool.Pool
	Config     *config.Config
	HTTPAddr   string
}

// Build initializes all database pools, repositories, use cases, and HTTP handlers for the merchant service.
func Build(cfg *config.Config) (*App, error) {
	ctx := context.Background()

	pool, err := postgres.ConnectDB(ctx, *cfg)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}

	migrationDSN := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=merchant",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)
	if err := postgres.RunMigrations(migrationDSN, cfg.Database.MigrationsPath); err != nil {
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	storeRepo := repo.NewStoreRepo(pool)
	storeUseCase := usecases.NewStoreUseCase(storeRepo)

	campaignRepo := repo.NewCampaignRepo(pool)
	campaignUseCase := usecases.NewCampaignUseCase(campaignRepo, storeRepo)

	rewardRepo := repo.NewRewardRepo(pool)
	rewardUseCase := usecases.NewRewardUseCase(rewardRepo, storeRepo)

	subscriptionRepo := repo.NewSubscriptionRepo(pool)
	subscriptionUseCase := usecases.NewSubscriptionUseCase(subscriptionRepo, storeRepo, campaignRepo)

	kycRepo := repo.NewKYCRepo(pool)
	kycUseCase := usecases.NewKYCUseCase(kycRepo)

	qrRepo := repo.NewQRRepo(pool)

	dashboardUseCase := usecases.NewDashboardUseCase(storeRepo, rewardRepo, campaignRepo, qrRepo)

	adminRepo := repo.NewAdminRepo(pool, storeRepo, campaignRepo, kycRepo)
	adminUseCase := usecases.NewAdminUseCase(adminRepo)

	// Websocket Hub Setup
	hub := ws.NewHub()
	hub.StartCleanupRoutine(5 * time.Minute)

	// HTTP Setup
	httpHandler := transporthttp.NewMerchantHandler(storeUseCase)
	campaignHandler := transporthttp.NewCampaignHandler(campaignUseCase)
	rewardHandler := transporthttp.NewRewardHandler(rewardUseCase)
	subscriptionHandler := transporthttp.NewSubscriptionHandler(subscriptionUseCase)
	kycHandler := transporthttp.NewKYCHandler(kycUseCase)
	dashboardHandler := transporthttp.NewDashboardHandler(dashboardUseCase)
	adminHandler := transporthttp.NewAdminHandler(adminUseCase)
	httpServer := setupHTTPServer(httpHandler, campaignHandler, rewardHandler, subscriptionHandler, kycHandler, dashboardHandler, adminHandler, hub)

	httpAddr := cfg.Server.HTTPAddr
	if port := os.Getenv("PORT"); port != "" {
		httpAddr = port
	}
	if httpAddr == "" {
		httpAddr = "9093"
	}
	if !strings.Contains(httpAddr, ":") {
		httpAddr = "0.0.0.0:" + httpAddr
	}

	return &App{
		HTTPServer: httpServer,
		DB:         pool,
		Config:     cfg,
		HTTPAddr:   httpAddr,
	}, nil
}

// setupHTTPServer configures Fiber middleware, CORS policies, health check, and registers all REST routes.
func setupHTTPServer(handler *transporthttp.MerchantHandler, campaignHandler *transporthttp.CampaignHandler, rewardHandler *transporthttp.RewardHandler, subscriptionHandler *transporthttp.SubscriptionHandler, kycHandler *transporthttp.KYCHandler, dashboardHandler *transporthttp.DashboardHandler, adminHandler *transporthttp.AdminHandler, hub *ws.Hub) *fiber.App {
	server := fiber.New(fiber.Config{
		AppName: "merchant-service",
	})

	server.Use(recover.New())
	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	server.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": "merchant"})
	})

	// No-op auth middleware for now, or use real one when JWT config is wired
	noopAuth := func(c fiber.Ctx) error { return c.Next() }
	transporthttp.RegisterRoutes(server, handler, campaignHandler, rewardHandler, subscriptionHandler, kycHandler, dashboardHandler, noopAuth)
	transporthttp.RegisterAdminRoutes(server, adminHandler, noopAuth)
	transporthttp.RegisterWebSocketRoutes(server, hub.HandleConnection())

	return server
}




func (a *App) Run() error {
	// Start HTTP server
	log.Printf("Starting HTTP server on %s", a.HTTPAddr)
	return a.HTTPServer.Listen(a.HTTPAddr)
}

func (a *App) Close() {
	if a.HTTPServer != nil {
		a.HTTPServer.Shutdown()
	}
	if a.DB != nil {
		a.DB.Close()
	}
}
