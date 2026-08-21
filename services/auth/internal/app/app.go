package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moneymate-2026/moneymate-backend/auth/config"
	"github.com/moneymate-2026/moneymate-backend/auth/internal/adapter/postgres"
	"github.com/moneymate-2026/moneymate-backend/auth/internal/adapter/postgres/repo"
	rediscard "github.com/moneymate-2026/moneymate-backend/auth/internal/adapter/redis"
	"github.com/moneymate-2026/moneymate-backend/auth/internal/domain"
	"github.com/moneymate-2026/moneymate-backend/auth/internal/infra/hasher"
	"github.com/moneymate-2026/moneymate-backend/auth/internal/infra/idgen"
	"github.com/moneymate-2026/moneymate-backend/auth/internal/infra/mailer"
	"github.com/moneymate-2026/moneymate-backend/auth/internal/infra/outboxpublisher"
	"github.com/moneymate-2026/moneymate-backend/auth/internal/infra/tokenissuer"
	transporthttp "github.com/moneymate-2026/moneymate-backend/auth/internal/transport/http"
	usecase "github.com/moneymate-2026/moneymate-backend/auth/internal/usecases"
	s3util "github.com/moneymate-2026/moneymate-backend/shared/pkg/S3"
	sharedjwt "github.com/moneymate-2026/moneymate-backend/shared/pkg/jwt"
	"github.com/moneymate-2026/moneymate-backend/shared/pkg/kafka"
	sharedmailer "github.com/moneymate-2026/moneymate-backend/shared/pkg/mailer"
	sharedpgxtx "github.com/moneymate-2026/moneymate-backend/shared/pkg/pgxtx"
	"github.com/redis/go-redis/v9"
)

const dbConnectTimeout = 10 * time.Second

type App struct {
	Server      *fiber.App
	DB          *pgxpool.Pool
	RedisClient *redis.Client
	Config      *config.Config
	Publisher   *outboxpublisher.Publisher // NEW
}

func Build(cfg *config.Config) (app *App, err error) {
	var pool *pgxpool.Pool
	var redisClient *redis.Client

	defer func() {
		if err != nil {
			if redisClient != nil {
				redisClient.Close()
			}
			if pool != nil {
				pool.Close()
			}
		}
	}()

	pool, err = connectDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s&search_path=auth",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SslMode,
	)
	if err = postgres.RunMigrations(dsn, cfg.Database.MigrationsPath); err != nil {
		return nil, fmt.Errorf("run migrations: %w", err)
	}
	staffRepo := repo.NewStaffRepo(pool)
	roleRepo := repo.NewRoleRepo(pool)
	if err := seedAdmin(context.Background(), staffRepo, roleRepo, hasher.New(), idgen.New(), cfg); err != nil {
		return nil, fmt.Errorf("seed admin: %w", err)
	}
	redisClient, err = setupRedis(cfg)
	if err != nil {
		return nil, fmt.Errorf("setup redis: %w", err)
	}

	outboxRepo := repo.NewOutboxRepo(pool)

	kafkaProducer, err := kafka.NewProducer(kafka.Config{
		Brokers:  cfg.Kafka.Brokers,
		Username: cfg.Kafka.Username,
		Password: cfg.Kafka.Password,
		CACert:   cfg.Kafka.CACert,
	})
	if err != nil {
		return nil, fmt.Errorf("create kafka producer: %w", err)
	}
	s3Client, err := s3util.New(context.Background(), cfg.S3.Bucket, cfg.S3.Region, cfg.S3.PublicBase)
	if err != nil {
		return nil, fmt.Errorf("init s3 client: %w", err)
	}

	publisher := outboxpublisher.New(outboxRepo, kafkaProducer)

	handlers := setupDependencies(pool, redisClient, cfg, outboxRepo, s3Client)
	server := setupServer(cfg, handlers, pool, redisClient)

	return &App{
		Server:      server,
		DB:          pool,
		RedisClient: redisClient,
		Config:      cfg,
		Publisher:   publisher,
	}, nil
}

func (a *App) Close() {
	if a.DB != nil {
		a.DB.Close()
	}
	if a.RedisClient != nil {
		a.RedisClient.Close()
	}
}

// ─── Private Setup Helpers ───────────────────────────────────────────────

func connectDB(cfg *config.Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbConnectTimeout)
	defer cancel()
	return postgres.ConnectDB(ctx, cfg.Database.DSN)
}

func setupRedis(cfg *config.Config) (*redis.Client, error) {
	return rediscard.NewClient(rediscard.Config{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       0,
	})
}

func setupDependencies(pool *pgxpool.Pool, redisClient *redis.Client, cfg *config.Config, outboxRepo domain.OutboxRepository, s3Client *s3util.Client) *transporthttp.Handlers {
	jwtCfg := sharedjwt.Config{
		AccessSecret:      cfg.JWT.AccessSecret,
		RefreshSecret:     cfg.JWT.RefreshSecret,
		AccessExpiryMins:  cfg.JWT.AccessExpiryMinutes,
		RefreshExpiryHrs:  cfg.JWT.RefreshExpiryHours,
		TxTokenExpirySecs: cfg.JWT.TxTokenExpirySecs,
	}

	emailCfg := sharedmailer.Config{
		APIKey:      cfg.Email.APIKey,
		FromAddress: cfg.Email.FromAddress,
		FromName:    cfg.Email.FromName,
	}

	h := hasher.New()
	g := idgen.New()
	issuer := tokenissuer.New(jwtCfg)
	mailerClient := sharedmailer.New(emailCfg)
	otpMailer := mailer.NewOtpMail(mailerClient)

	
	userRepo := repo.NewUserRepo(pool)
	staffRepo := repo.NewStaffRepo(pool)
	roleRepo := repo.NewRoleRepo(pool)
	refreshTokenRepo := repo.NewRefreshTokenRepo(pool)
	pinRepo := repo.NewUserPinRepo(pool)
	permRepo := repo.NewPermissionRepo(pool)
	store := rediscard.NewStore(redisClient)
	txMgr := sharedpgxtx.New(pool)

	pinUC := usecase.NewUserPinUsecase(pinRepo, h, g)
	authUC := usecase.NewAuthUsecase(userRepo, roleRepo, outboxRepo, refreshTokenRepo, pinRepo, pinUC, store, txMgr, h, g, issuer, jwtCfg, staffRepo)

	otpMailerIface := usecase.EmailSender(otpMailer)
	if cfg.Env == "dev" {
		otpMailerIface = mailer.NewDevOtpMail()
		log.Println("[DEV MODE] OTP codes will be logged to console instead of sent via email")
	}

	otpUC := usecase.NewOTPUsecase(userRepo, store, otpMailerIface, cfg.OTP)
	adminRoleUC := usecase.NewAdminRoleUsecase(roleRepo, userRepo, g)
	adminUserUC := usecase.NewAdminUserUsecase(userRepo, roleRepo, h, g, pinRepo, txMgr)
	staffUC := usecase.NewStaffUsecase(staffRepo, roleRepo, h, g)
	permissionUC := usecase.NewPermissionUsecase(permRepo, roleRepo, g)
	profilePictureUC := usecase.NewProfilePictureUsecase(userRepo, s3Client)

	return &transporthttp.Handlers{
		Auth:           transporthttp.NewAuthHandler(authUC, otpUC, userRepo, cfg.JWT.AccessSecret, redisClient),
		Role:           transporthttp.NewRoleHandler(adminRoleUC),
		User:           transporthttp.NewUserHandler(adminUserUC),
		Staff:          transporthttp.NewStaffHandler(staffUC),
		UserPin:        transporthttp.NewUserPinHandler(pinUC, issuer),
		Permission:     transporthttp.NewPermissionHandler(permissionUC),
		Profile: transporthttp.NewProfilePictureHandler(profilePictureUC),
	}
}

func setupServer(cfg *config.Config, handlers *transporthttp.Handlers, pool *pgxpool.Pool, redisClient *redis.Client) *fiber.App {
	server := fiber.New(fiber.Config{
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		AppName:      "auth-service",
	})

	server.Use(recover.New())
	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Device-Id"},
	}))

	registerHealthRoutes(server, pool, redisClient)
	transporthttp.RegisterRoutes(server, handlers, cfg.InternalServiceSecret)

	return server
}

func registerHealthRoutes(server *fiber.App, pool *pgxpool.Pool, redisClient *redis.Client) {
	server.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": "auth"})
	})

	server.Get("/ready", func(c fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "unavailable", "dependency": "postgres", "error": err.Error(),
			})
		}
		if err := redisClient.Ping(ctx).Err(); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "unavailable", "dependency": "redis", "error": err.Error(),
			})
		}
		return c.JSON(fiber.Map{"status": "ready", "service": "auth"})
	})
}
