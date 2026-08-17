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

	"github.com/moneymate-2026/moneymate-backend/services/notification/config"
	"github.com/moneymate-2026/moneymate-backend/services/notification/internal/adapter/fcm"
	"github.com/moneymate-2026/moneymate-backend/services/notification/internal/adapter/postgres"
	"github.com/moneymate-2026/moneymate-backend/services/notification/internal/adapter/postgres/repo"
	transporthttp "github.com/moneymate-2026/moneymate-backend/services/notification/internal/transport/http"
	"github.com/moneymate-2026/moneymate-backend/services/notification/internal/usecases"
	sharedjwt "github.com/moneymate-2026/moneymate-backend/shared/pkg/jwt"
	"github.com/moneymate-2026/moneymate-backend/shared/pkg/kafka"
)

var notificationTopics = []string{
	"moneymate.bills.due",
	"moneymate.debt.created",
	"moneymate.debt.confirmed",
	"moneymate.wallet.transfer.completed",
	"moneymate.wallet.topup.completed",
	"moneymate.payment.merchant.completed",
	"moneymate.merchant.commission.credited",
	"moneymate.merchant.verified",
	"moneymate.campaign.created",
	"moneymate.offer.redeemed",
}

type App struct {
	HTTPServer     *fiber.App
	DB             *pgxpool.Pool
	HTTPAddr       string
	KafkaConsumer  *kafka.Consumer
	NotificationUC *usecases.NotificationUsecase
	DeviceUC       *usecases.DeviceUsecase
	InboxUC        *usecases.InboxUsecase
	PreferenceUC   *usecases.PreferenceUsecase
}

func Build(cfg *config.Config) (*App, error) {
	ctx := context.Background()

	pool, err := postgres.ConnectDB(ctx, *cfg)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s&search_path=notification",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SslMode,
	)
	if err := postgres.RunMigrations(dsn, cfg.Database.MigrationsPath); err != nil {
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	deviceRepo := repo.NewDeviceRepo(pool)
	inboxRepo := repo.NewInboxRepo(pool)
	preferenceRepo := repo.NewPreferenceRepo(pool)
	deliveryRepo := repo.NewDeliveryRepo(pool)

	fcmClient, err := fcm.New(ctx, cfg.FCM.ProjectID, cfg.FCM.CredentialsPath)
	if err != nil {
		return nil, fmt.Errorf("init fcm: %w", err)
	}

	notificationUC := usecases.NewNotificationUsecase(inboxRepo, preferenceRepo, deviceRepo, deliveryRepo, fcmClient)
	deviceUC := usecases.NewDeviceUsecase(deviceRepo)
	inboxUC := usecases.NewInboxUsecase(inboxRepo)
	preferenceUC := usecases.NewPreferenceUsecase(preferenceRepo)

	jwtCfg := sharedjwt.Config{
		AccessSecret:     cfg.JWT.AccessSecret,
		RefreshSecret:    cfg.JWT.RefreshSecret,
		AccessExpiryMins: cfg.JWT.AccessExpiryMinutes,
		RefreshExpiryHrs: cfg.JWT.RefreshExpiryHours,
	}

	server := setupHTTPServer(deviceUC, inboxUC, preferenceUC, jwtCfg)

	kafkaConsumer, err := kafka.NewConsumer(kafka.ConsumerConfig{
		Brokers:     cfg.Kafka.Brokers,
		Username:    cfg.Kafka.Username,
		Password:    cfg.Kafka.Password,
		CACert:      cfg.Kafka.CACert,
		GroupTopics: notificationTopics,
		GroupID:     "moneymate-notification-svc",
	})
	if err != nil {
		return nil, fmt.Errorf("create kafka consumer: %w", err)
	}

	httpAddr := cfg.Server.HTTPAddr
	if port := os.Getenv("PORT"); port != "" {
		httpAddr = port
	}
	if httpAddr == "" {
		httpAddr = "9095"
	}
	if !strings.Contains(httpAddr, ":") {
		httpAddr = "0.0.0.0:" + httpAddr
	}

	return &App{
		HTTPServer:     server,
		DB:             pool,
		HTTPAddr:       httpAddr,
		KafkaConsumer:  kafkaConsumer,
		NotificationUC: notificationUC,
		DeviceUC:       deviceUC,
		InboxUC:        inboxUC,
		PreferenceUC:   preferenceUC,
	}, nil
}

func setupHTTPServer(du *usecases.DeviceUsecase, iu *usecases.InboxUsecase, pu *usecases.PreferenceUsecase, jwtCfg sharedjwt.Config) *fiber.App {
	server := fiber.New(fiber.Config{AppName: "notification-service"})
	server.Use(recover.New())
	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	server.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": "notification"})
	})

	transporthttp.RegisterRoutes(server, du, iu, pu, jwtCfg)
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
	if a.KafkaConsumer != nil {
		a.KafkaConsumer.Close()
	}
}
