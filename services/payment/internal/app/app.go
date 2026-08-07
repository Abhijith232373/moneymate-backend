package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moneymate-2026/moneymate-backend/services/payment/config"
	"github.com/moneymate-2026/moneymate-backend/services/payment/internal/adapter/postgres"
	"github.com/moneymate-2026/moneymate-backend/services/payment/internal/adapter/postgres/repo"
	transporthttp "github.com/moneymate-2026/moneymate-backend/services/payment/internal/transport/http"
	"github.com/moneymate-2026/moneymate-backend/services/payment/internal/usecases"
	sharedjwt "github.com/moneymate-2026/moneymate-backend/shared/pkg/jwt"
)

type App struct {
	HTTPServer *fiber.App
	DB         *pgxpool.Pool
	HTTPAddr   string
}

func Build(cfg *config.Config) (*App, error) {
	ctx := context.Background()

	pool, err := postgres.ConnectDB(ctx, *cfg)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}

	accountRepo := repo.NewAccountRepo(pool)
	transactionRepo := repo.NewTransactionRepo(pool)
	ledgerRepo := repo.NewLedgerRepo(pool)

	walletUC := usecases.NewWalletUsecase(accountRepo)
	transferUC := usecases.NewTransferUsecase(accountRepo, transactionRepo, ledgerRepo)

	walletHandler := transporthttp.NewWalletHandler(walletUC)
	transferHandler := transporthttp.NewTransferHandler(transferUC)

	jwtCfg := sharedjwt.Config{
		AccessSecret:     cfg.JWT.AccessSecret,
		RefreshSecret:    cfg.JWT.RefreshSecret,
		AccessExpiryMins: cfg.JWT.AccessExpiryMinutes,
		RefreshExpiryHrs: cfg.JWT.RefreshExpiryHours,
	}

	server := setupHTTPServer(walletHandler, transferHandler, jwtCfg)

	httpAddr := cfg.Server.HTTPAddr
	if port := os.Getenv("PORT"); port != "" {
		httpAddr = port
	}
	if httpAddr == "" {
		httpAddr = "9094"
	}
	if !strings.Contains(httpAddr, ":") {
		httpAddr = "0.0.0.0:" + httpAddr
	}

	return &App{HTTPServer: server, DB: pool, HTTPAddr: httpAddr}, nil
}

func setupHTTPServer(wh *transporthttp.WalletHandler, th *transporthttp.TransferHandler, jwtCfg sharedjwt.Config) *fiber.App {
	server := fiber.New(fiber.Config{AppName: "payment-service"})
	server.Use(recover.New())
	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Transaction-Token"},
	}))

	server.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": "payment"})
	})

	transporthttp.RegisterRoutes(server, wh, th, jwtCfg)
	return server
}

func (a *App) Run() error {
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
