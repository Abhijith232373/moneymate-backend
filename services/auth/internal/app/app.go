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
	"github.com/redis/go-redis/v9"

	"github.com/moneymate-2026/moneymate-backend/auth/config"
	"github.com/moneymate-2026/moneymate-backend/auth/internal/adapter/postgres"
	"github.com/moneymate-2026/moneymate-backend/auth/internal/adapter/postgres/repo"
	rediscard "github.com/moneymate-2026/moneymate-backend/auth/internal/adapter/redis"
	"github.com/moneymate-2026/moneymate-backend/auth/internal/infra/hasher"
	"github.com/moneymate-2026/moneymate-backend/auth/internal/infra/idgen"
	"github.com/moneymate-2026/moneymate-backend/auth/internal/infra/mailer"
	"github.com/moneymate-2026/moneymate-backend/auth/internal/infra/tokenissuer"
	transporthttp "github.com/moneymate-2026/moneymate-backend/auth/internal/transport/http"
	usecase "github.com/moneymate-2026/moneymate-backend/auth/internal/usecases"
	sharedjwt "github.com/moneymate-2026/moneymate-backend/shared/pkg/jwt"
	sharedmailer "github.com/moneymate-2026/moneymate-backend/shared/pkg/mailer"
	sharedpgxtx "github.com/moneymate-2026/moneymate-backend/shared/pkg/pgxtx"
)

const dbConnectTimeout = 10 * time.Second

type App struct {
	Server      *fiber.App
	DB          *pgxpool.Pool
	RedisClient *redis.Client
	Config      *config.Config
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
	userRepo := repo.NewUserRepo(pool)
	roleRepo := repo.NewRoleRepo(pool)
	if err := seedAdmin(context.Background(), userRepo, roleRepo, hasher.New(), idgen.New(), cfg); err != nil {
		return nil, fmt.Errorf("seed admin: %w", err)
	}
	redisClient, err = setupRedis(cfg)
	if err != nil {
		return nil, fmt.Errorf("setup redis: %w", err)
	}

	handlers := setupDependencies(pool, redisClient, cfg)
	server := setupServer(cfg, handlers, pool, redisClient)

	return &App{
		Server:      server,
		DB:          pool,
		RedisClient: redisClient,
		Config:      cfg,
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

func setupDependencies(pool *pgxpool.Pool, redisClient *redis.Client, cfg *config.Config) *transporthttp.Handlers {
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

	// Infra
	h := hasher.New()
	g := idgen.New()
	issuer := tokenissuer.New(jwtCfg)
	mailerClient := sharedmailer.New(emailCfg)
	otpMailer := mailer.NewOtpMail(mailerClient)

	// Repositories
	userRepo := repo.NewUserRepo(pool)
	roleRepo := repo.NewRoleRepo(pool)
	refreshTokenRepo := repo.NewRefreshTokenRepo(pool)
	pinRepo := repo.NewUserPinRepo(pool)
	permRepo := repo.NewPermissionRepo(pool) // NEW
	store := rediscard.NewStore(redisClient)
	txMgr := sharedpgxtx.New(pool)

	// Usecases
	pinUC := usecase.NewUserPinUsecase(pinRepo, h, g)
	authUC := usecase.NewAuthUsecase(userRepo, roleRepo, refreshTokenRepo, pinRepo, pinUC, store, txMgr, h, g, issuer, jwtCfg)

	otpMailerIface := usecase.EmailSender(otpMailer)
	if cfg.Env == "dev" {
		otpMailerIface = mailer.NewDevOtpMail()
		log.Println("[DEV MODE] OTP codes will be logged to console instead of sent via email")
	}

	otpUC := usecase.NewOTPUsecase(userRepo, store, otpMailerIface, cfg.OTP)
	adminRoleUC := usecase.NewAdminRoleUsecase(roleRepo, userRepo, g)
	adminUserUC := usecase.NewAdminUserUsecase(userRepo, roleRepo, h, g)
	permissionUC := usecase.NewPermissionUsecase(permRepo, roleRepo, g) // NEW

	return &transporthttp.Handlers{
		Auth:       transporthttp.NewAuthHandler(authUC, otpUC, userRepo, cfg.JWT.AccessSecret, redisClient),
		Role:       transporthttp.NewRoleHandler(adminRoleUC),
		User:       transporthttp.NewUserHandler(adminUserUC),
		UserPin:    transporthttp.NewUserPinHandler(pinUC, issuer),
		Permission: transporthttp.NewPermissionHandler(permissionUC), // NEW
	}
}

// setupServer configures the Fiber application, middleware, and routing.
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