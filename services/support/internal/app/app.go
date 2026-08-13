package app

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	_ "github.com/lib/pq"
	"github.com/abijith/moneymate-backend/services/support/config"
	"github.com/abijith/moneymate-backend/services/support/internal/repository/postgres"
	"github.com/abijith/moneymate-backend/services/support/internal/usecase"
	transporthttp "github.com/abijith/moneymate-backend/services/support/internal/transport/http"
)

type App struct {
	HTTPServer *fiber.App
	DB         *sql.DB
	Config     *config.Config
	HTTPAddr   string
}

func Build(cfg *config.Config) (*App, error) {
	dbURL := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SslMode,
	)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	supportRepo := postgres.NewSupportRepo(db)
	supportUseCase := usecase.NewSupportUseCase(supportRepo)

	supportHandler := transporthttp.NewSupportHandler(supportUseCase)
	httpServer := setupHTTPServer(supportHandler)

	httpAddr := cfg.Server.HTTPAddr
	if port := os.Getenv("PORT"); port != "" {
		httpAddr = port
	}
	if httpAddr == "" {
		httpAddr = "8085"
	}
	if !strings.Contains(httpAddr, ":") {
		httpAddr = "0.0.0.0:" + httpAddr
	}

	return &App{
		HTTPServer: httpServer,
		DB:         db,
		Config:     cfg,
		HTTPAddr:   httpAddr,
	}, nil
}

func setupHTTPServer(handler *transporthttp.SupportHandler) *fiber.App {
	server := fiber.New(fiber.Config{
		AppName: "support-service",
	})

	server.Use(recover.New())
	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	noopAuth := func(c fiber.Ctx) error { return c.Next() }

	transporthttp.RegisterRoutes(server, handler, noopAuth)
	transporthttp.RegisterAdminRoutes(server, handler, noopAuth)

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
