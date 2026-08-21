# Reward Service Implementation Guide

This guide maps the current MoneyMate backend and gives a concrete route plan for adding an independent rewards service without coupling it to payment internals too early.

The recommendation is to build the new service as `services/rewards`, using the already bootstrapped Postgres schema/user named `rewards` and `rewards_user`.

## Current Backend Map

The backend is a Go workspace at the root:

```text
moneymate-backend/
  go.work
  docker-compose.yml
  docker-compose.prod.yml
  Taskfile.yml
  infra/postgres/bootstrap/
  shared/
  services/
    auth/
    gateway/
    merchant/
    notification/
    payment/
    support/
```

Important existing patterns:

```text
auth
  Own schema: auth
  Startup: config.LoadConfig -> app.Build -> Fiber server
  Uses: Postgres, Redis, Kafka producer/outbox
  Health: /health and /ready with DB/Redis checks
  Internal auth: /internal/... protected by X-Internal-Secret

payment
  Own schema: payment
  Startup: config.LoadConfig -> app.Build -> Fiber server + Kafka consumer
  Uses: Postgres, Kafka consumer, shared JWT, internal auth client
  Routes: /payment/... for user APIs, /internal/... protected by X-Internal-Secret
  Money unit: int64 paise

notification
  Own schema: notification
  Startup: Fiber server + Kafka group consumer
  Good reference for consuming multiple event topics

merchant
  Own schema: merchant
  Has an existing Rewards Center:
    /merchant/rewards/summary
    /merchant/rewards/history
    /merchant/rewards/redeem
  This is merchant/store loyalty balance UI logic, not the independent cashback reward processor.

gateway
  Client entry point under /api/v1
  Authenticates with shared JWT middleware
  Proxies to downstream services via proxy.HTTPProxy and ServiceRegistry
  Current registry only includes payment and support
```

Database bootstrap already includes:

```sql
CREATE SCHEMA IF NOT EXISTS rewards;
CREATE ROLE rewards_user LOGIN PASSWORD 'rewards_password';
GRANT USAGE, CREATE ON SCHEMA rewards TO rewards_user;
GRANT CREATE ON DATABASE moneymate TO rewards_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA rewards
  GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO rewards_user;
```

So the new service should use schema name `rewards`, env vars `REWARDS_DB_USER` and `REWARDS_DB_PASSWORD`, and service key `rewards`.

## Route Map

### Direct rewards service routes

These are the routes served by `services/rewards` itself.

```text
GET    /health
GET    /ready

GET    /rewards/me
GET    /rewards

POST   /admin/rewards/rules
GET    /admin/rewards/rules
GET    /admin/rewards/rules/:id
PUT    /admin/rewards/rules/:id
PATCH  /admin/rewards/rules/:id/deactivate

GET    /internal/rewards/payouts/:id
POST   /internal/rewards/replay-failed
```

Recommended behavior:

```text
GET /rewards/me
  Auth: user JWT
  Returns reward payouts where recipient_id equals authenticated user_id.
  Query params: limit, offset, status

GET /rewards?transaction_id=...
  Auth: user JWT
  Returns reward payouts tied to one original payment transaction.
  Must verify the result belongs to the current user unless admin/internal.

POST /admin/rewards/rules
  Auth: gateway admin role
  Creates a reward rule.

GET /admin/rewards/rules
  Auth: gateway admin role
  Lists all rules, newest first.

PUT /admin/rewards/rules/:id
  Auth: gateway admin role
  Updates rule config.

PATCH /admin/rewards/rules/:id/deactivate
  Auth: gateway admin role
  Sets active=false.
```

### Gateway routes

Add a dedicated gateway route file, matching `payment_routes.go`:

```text
GET    /api/v1/rewards/me
GET    /api/v1/rewards?transaction_id=...

POST   /api/v1/admin/rewards/rules
GET    /api/v1/admin/rewards/rules
GET    /api/v1/admin/rewards/rules/:id
PUT    /api/v1/admin/rewards/rules/:id
PATCH  /api/v1/admin/rewards/rules/:id/deactivate
```

Gateway route mapping:

```go
func registerRewardRoutes(api fiber.Router, authMiddleware fiber.Handler, registry *proxy.ServiceRegistry) {
	rewards := api.Group("/rewards")
	rewards.Use(authMiddleware)
	rewards.Get("/me", proxy.HTTPProxy(registry, "rewards", "/rewards/me"))
	rewards.Get("/", proxy.HTTPProxy(registry, "rewards", "/rewards"))
}

func registerAdminRewardRuleRoutes(admin fiber.Router, registry *proxy.ServiceRegistry) {
	rules := admin.Group("/rewards/rules")
	rules.Post("/", proxy.HTTPProxy(registry, "rewards", "/admin/rewards/rules"))
	rules.Get("/", proxy.HTTPProxy(registry, "rewards", "/admin/rewards/rules"))
	rules.Get("/:id", proxy.HTTPProxy(registry, "rewards", "/admin/rewards/rules/:id"))
	rules.Put("/:id", proxy.HTTPProxy(registry, "rewards", "/admin/rewards/rules/:id"))
	rules.Patch("/:id/deactivate", proxy.HTTPProxy(registry, "rewards", "/admin/rewards/rules/:id/deactivate"))
}
```

Also add `rewards` to:

```text
services/gateway/config/config.go
services/gateway/config/config.yaml
docker-compose.yml gateway environment
docker-compose.prod.yml gateway environment
```

## Implementation Order For Today

Follow this order. It avoids blocking on the payment service.

```text
1. Create service skeleton and boot it with /health and /ready.
2. Add migrations, SQLC config, and generated query layer.
3. Implement pure rule engine with table-driven tests.
4. Implement reward rule admin CRUD.
5. Add Kafka consumer skeleton and local fake-event publisher.
6. Add payout execution behind PaymentClient with a fake implementation.
7. Add user/frontend query endpoints.
8. Register gateway routes and docker/task wiring.
```

Only the real `PaymentClient` implementation is blocked until payment adds an internal payout endpoint.

## How To Use This Guide

All paths in this guide are relative to:

```text
C:\Users\Fazal\Money-Mate\moneymate-backend
```

Example:

```text
services/rewards/cmd/main.go
```

means:

```text
C:\Users\Fazal\Money-Mate\moneymate-backend\services\rewards\cmd\main.go
```

The guide uses this wording:

```text
Create file
  The file does not exist yet. Create it at the exact path shown.

Edit file
  The file already exists. Add or change only the lines mentioned.

Copy from
  Start by copying an existing file from another service, then change names/imports.

Run
  Command to run from C:\Users\Fazal\Money-Mate\moneymate-backend unless another directory is shown.
```

Important rule:

```text
Do not edit files inside sqlc/generated manually.
Those files are generated by sqlc.
```

## Detailed File-By-File Checklist

This section is the simple-English checklist. Do it in order.

### Phase 1: Create The Rewards Service Folder

Run these commands from:

```text
C:\Users\Fazal\Money-Mate\moneymate-backend
```

PowerShell commands:

```powershell
New-Item -ItemType Directory -Force services\rewards
New-Item -ItemType Directory -Force services\rewards\cmd
New-Item -ItemType Directory -Force services\rewards\config
New-Item -ItemType Directory -Force services\rewards\internal\app
New-Item -ItemType Directory -Force services\rewards\internal\adapter\postgres\repo
New-Item -ItemType Directory -Force services\rewards\internal\adapter\paymentclient
New-Item -ItemType Directory -Force services\rewards\internal\domain
New-Item -ItemType Directory -Force services\rewards\internal\infra\kafkaconsumer
New-Item -ItemType Directory -Force services\rewards\internal\transport\http
New-Item -ItemType Directory -Force services\rewards\internal\usecases
New-Item -ItemType Directory -Force services\rewards\migrations
New-Item -ItemType Directory -Force services\rewards\sqlc\queries
New-Item -ItemType Directory -Force services\rewards\sqlc\generated
```

Create these empty files next:

```powershell
New-Item -ItemType File -Force services\rewards\go.mod
New-Item -ItemType File -Force services\rewards\Dockerfile
New-Item -ItemType File -Force services\rewards\.air.toml
New-Item -ItemType File -Force services\rewards\cmd\main.go
New-Item -ItemType File -Force services\rewards\config\config.go
New-Item -ItemType File -Force services\rewards\config\config.yaml
New-Item -ItemType File -Force services\rewards\internal\app\app.go
New-Item -ItemType File -Force services\rewards\internal\adapter\postgres\db.go
New-Item -ItemType File -Force services\rewards\internal\transport\http\routes.go
New-Item -ItemType File -Force services\rewards\internal\transport\http\middleware.go
New-Item -ItemType File -Force services\rewards\sqlc\sqlc.yml
```

### Phase 2: Add The Module File

Create file:

```text
services/rewards/go.mod
```

Paste:

```go
module github.com/moneymate-2026/moneymate-backend/services/rewards

go 1.26.4

require (
	github.com/gofiber/fiber/v3 v3.4.0
	github.com/golang-migrate/migrate/v4 v4.19.1
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.10.0
	github.com/joho/godotenv v1.5.1
	github.com/moneymate-2026/moneymate-backend/shared v0.0.0
	github.com/spf13/viper v1.21.0
)
```

Then run:

```powershell
go work use .\services\rewards
```

Check:

```powershell
Get-Content go.work
```

You should see:

```text
./services/rewards
```

### Phase 3: Add Hot Reload Config

Copy from:

```text
services/payment/.air.toml
```

To:

```text
services/rewards/.air.toml
```

You usually do not need to change anything inside this file. It only tells Air how to rebuild the Go service while Docker is running.

### Phase 4: Add Dockerfile

Create file:

```text
services/rewards/Dockerfile
```

Use `services/payment/Dockerfile` as the template. Change only payment-specific names to rewards.

Use this final shape:

```dockerfile
FROM golang:1.26-alpine AS base
WORKDIR /moneymate-backend

RUN apk add --no-cache git
ENV GONOSUMCHECK=* GONOSUMDB=*
ENV GOFLAGS=-buildvcs=false

COPY go.work go.work.sum ./
COPY shared/go.mod shared/go.mod
COPY services/auth/go.mod services/auth/go.mod
COPY services/gateway/go.mod services/gateway/go.mod
COPY services/merchant/go.mod services/merchant/go.mod
COPY services/payment/go.mod services/payment/go.mod
COPY services/support/go.mod services/support/go.mod
COPY services/notification/go.mod services/notification/go.mod
COPY services/rewards/go.mod services/rewards/go.mod

RUN go work sync

FROM base AS dev
WORKDIR /moneymate-backend/services/rewards
RUN go install github.com/air-verse/air@latest
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
CMD ["air", "-c", ".air.toml"]

FROM base AS builder
WORKDIR /moneymate-backend
COPY shared/ shared/
COPY services/rewards/ services/rewards/
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/rewards ./services/rewards/cmd/main.go

FROM alpine:3.19 AS prod
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app

RUN mkdir -p /app/config
COPY --from=builder /bin/rewards ./rewards
COPY --from=builder /moneymate-backend/services/rewards/config/config.yaml ./config/config.yaml
COPY --from=builder /moneymate-backend/services/rewards/migrations ./migrations
COPY --from=builder /moneymate-backend/shared/certs/aiven-ca.pem ./certs/aiven-ca.pem

USER appuser
EXPOSE 9096
CMD ["./rewards"]
```

Important follow-up:

```text
After go.work contains ./services/rewards, update every existing service Dockerfile base stage.
```

Edit each file:

```text
services/auth/Dockerfile
services/gateway/Dockerfile
services/merchant/Dockerfile
services/payment/Dockerfile
services/support/Dockerfile
services/notification/Dockerfile
```

Add this line near the other `COPY services/.../go.mod` lines:

```dockerfile
COPY services/rewards/go.mod services/rewards/go.mod
```

Why:

```text
All Dockerfiles run go work sync.
If go.work mentions services/rewards but Docker did not copy services/rewards/go.mod,
Docker builds for other services can fail.
```

### Phase 5: Add Config Files

Create file:

```text
services/rewards/config/config.yaml
```

Paste:

```yaml
server:
  http_addr: "9096"

database:
  max_open_conns: 25
  min_open_conns: 5
  max_conn_lifetime: 15m
  max_idle_time: 5m
  migrations_path: "/app/migrations"

jwt:
  access_expiry_minutes: 15
  refresh_expiry_hours: 720
  tx_token_expiry_secs: 60

rewards:
  payment_completed_topic: "moneymate.payment.completed"
  consumer_group: "moneymate-rewards-svc"
  fake_payment_client: true
```

Create file:

```text
services/rewards/config/config.go
```

Paste the config loader from Step 1 below.

Plain-English explanation:

```text
LoadConfig reads services/rewards/config/config.yaml.
It also reads .env values like POSTGRES_HOST and REWARDS_DB_USER.
LoadDatabaseConfig(v, "rewards") tells Postgres to use the rewards schema.
LoadKafkaConfig reads Kafka credentials from .env.
LoadJWTConfig reads the shared JWT secret.
```

### Phase 6: Add Postgres Connection And Migrations Runner

Create file:

```text
services/rewards/internal/adapter/postgres/db.go
```

Copy from:

```text
services/payment/internal/adapter/postgres/db.go
```

Then change the import:

```go
"github.com/moneymate-2026/moneymate-backend/services/payment/config"
```

to:

```go
"github.com/moneymate-2026/moneymate-backend/services/rewards/config"
```

No other logic needs to change. This file should:

```text
Parse the DB DSN.
Create a pgx connection pool.
Ping Postgres.
Run migrations from /app/migrations.
```

### Phase 7: Add Main Entry Point

Create file:

```text
services/rewards/cmd/main.go
```

Use this simple version first:

```go
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
```

### Phase 8: Add App Bootstrap

Create file:

```text
services/rewards/internal/app/app.go
```

This file should do the same job as payment's `internal/app/app.go`, but simpler.

It should contain:

```text
type App struct
Build(cfg)
setupHTTPServer(...)
Run()
Close()
HandleKafkaMessage(...)
```

Use this starting version:

```go
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

	"github.com/moneymate-2026/moneymate-backend/services/rewards/config"
	"github.com/moneymate-2026/moneymate-backend/services/rewards/internal/adapter/postgres"
	transporthttp "github.com/moneymate-2026/moneymate-backend/services/rewards/internal/transport/http"
	"github.com/moneymate-2026/moneymate-backend/shared/pkg/kafka"
)

type App struct {
	HTTPServer    *fiber.App
	DB            *pgxpool.Pool
	HTTPAddr      string
	KafkaConsumer *kafka.Consumer
}

func Build(cfg *config.Config) (*App, error) {
	ctx := context.Background()

	pool, err := postgres.ConnectDB(ctx, *cfg)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s&search_path=rewards",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SslMode,
	)
	if err := postgres.RunMigrations(dsn, cfg.Database.MigrationsPath); err != nil {
		pool.Close()
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	kafkaConsumer, err := kafka.NewConsumer(kafka.ConsumerConfig{
		Brokers:  cfg.Kafka.Brokers,
		Username: cfg.Kafka.Username,
		Password: cfg.Kafka.Password,
		CACert:   cfg.Kafka.CACert,
		Topic:    cfg.Rewards.PaymentCompletedTopic,
		GroupID:  cfg.Rewards.ConsumerGroup,
	})
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("create kafka consumer: %w", err)
	}

	server := setupHTTPServer(pool, kafkaConsumer)

	httpAddr := cfg.Server.HTTPAddr
	if port := os.Getenv("PORT"); port != "" {
		httpAddr = port
	}
	if httpAddr == "" {
		httpAddr = "9096"
	}
	if !strings.Contains(httpAddr, ":") {
		httpAddr = "0.0.0.0:" + httpAddr
	}

	return &App{
		HTTPServer:    server,
		DB:            pool,
		HTTPAddr:      httpAddr,
		KafkaConsumer: kafkaConsumer,
	}, nil
}

func setupHTTPServer(pool *pgxpool.Pool, kafkaConsumer *kafka.Consumer) *fiber.App {
	server := fiber.New(fiber.Config{AppName: "rewards-service"})

	server.Use(recover.New())
	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Internal-Secret"},
	}))

	server.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": "rewards"})
	})

	server.Get("/ready", func(c fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "unavailable",
				"dependency": "postgres",
				"error": err.Error(),
			})
		}
		if kafkaConsumer == nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "unavailable",
				"dependency": "kafka",
				"error": "consumer not initialized",
			})
		}
		return c.JSON(fiber.Map{"status": "ready", "service": "rewards"})
	})

	transporthttp.RegisterRoutes(server)
	return server
}

func (a *App) Run() error {
	log.Printf("starting rewards HTTP server on %s", a.HTTPAddr)
	return a.HTTPServer.Listen(a.HTTPAddr)
}

func (a *App) Close() {
	if a.HTTPServer != nil {
		_ = a.HTTPServer.Shutdown()
	}
	if a.DB != nil {
		a.DB.Close()
	}
	if a.KafkaConsumer != nil {
		_ = a.KafkaConsumer.Close()
	}
}

func (a *App) HandleKafkaMessage(ctx context.Context, payload []byte) error {
	log.Printf("rewards kafka message received: %d bytes", len(payload))
	return nil
}
```

This is enough to boot the service before business logic exists.

### Phase 9: Add Empty Route File

Create file:

```text
services/rewards/internal/transport/http/routes.go
```

Paste:

```go
package http

import "github.com/gofiber/fiber/v3"

func RegisterRoutes(router fiber.Router) {
	rewards := router.Group("/rewards")

	rewards.Get("/me", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true, "data": []any{}})
	})

	rewards.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true, "data": []any{}})
	})
}
```

Later you will replace the placeholder handlers with real handler methods.

### Phase 10: Add Internal Secret Middleware

Create file:

```text
services/rewards/internal/transport/http/middleware.go
```

Paste:

```go
package http

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	sharedjwt "github.com/moneymate-2026/moneymate-backend/shared/pkg/jwt"
	response "github.com/moneymate-2026/moneymate-backend/shared/pkg/responses"
)

const localsUserID = "userID"

func RequireInternalSecret(secret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		provided := c.Get("X-Internal-Secret")
		if provided == "" || provided != secret {
			return response.Unauthorized(c, "internal access required")
		}
		return c.Next()
	}
}

func RequireUserID(cfg sharedjwt.Config) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return response.Unauthorized(c, "authentication required")
		}

		claims, err := sharedjwt.ParseAccessToken(parts[1], cfg.AccessSecret)
		if err != nil {
			return response.Unauthorized(c, "authentication required")
		}

		c.Locals(localsUserID, claims.UserID)
		return c.Next()
	}
}

func userIDFromLocals(c fiber.Ctx) string {
	id, _ := c.Locals(localsUserID).(string)
	return id
}
```

Use `RequireUserID` for `/rewards/...` service routes.
Use `RequireInternalSecret` for `/internal/...` routes.
For admin routes, rely on the gateway admin middleware first, then later add direct JWT role validation if needed.

### Phase 11: Add Migrations

Create file:

```text
services/rewards/migrations/000001_create_rewards_tables.up.sql
```

Paste the up migration from Step 3.

Create file:

```text
services/rewards/migrations/000001_create_rewards_tables.down.sql
```

Paste the down migration from Step 3.

Then run:

```powershell
task migrate:rewards
```

If the Taskfile task does not exist yet, run after adding the Taskfile changes in Phase 15.

### Phase 12: Add SQLC Files

Create file:

```text
services/rewards/sqlc/sqlc.yml
```

Paste the SQLC config from Step 4.

Create file:

```text
services/rewards/sqlc/queries/reward_rules.sql
```

Paste the `reward_rules.sql` code from Step 4.

Create file:

```text
services/rewards/sqlc/queries/reward_payouts.sql
```

Paste the `reward_payouts.sql` code from Step 4.

Then run:

```powershell
task sqlc:rewards
```

Expected result:

```text
SQLC creates Go files inside services/rewards/sqlc/generated.
Do not edit those generated files by hand.
```

### Phase 13: Add Domain Files

Create file:

```text
services/rewards/internal/domain/reward.go
```

Paste the `reward.go` domain code from Step 5.

Create file:

```text
services/rewards/internal/domain/payment.go
```

Paste the `payment.go` interface code from Step 5.

Plain-English explanation:

```text
domain/reward.go describes the reward data and the database methods the app needs.
domain/payment.go describes the payment service call we need later.
The domain package should not import Fiber, SQLC, or Kafka.
```

### Phase 14: Add Rule Engine And Tests

Create file:

```text
services/rewards/internal/usecases/rule_engine.go
```

Paste the rule engine code from Step 6.

Create file:

```text
services/rewards/internal/usecases/rule_engine_test.go
```

Paste the table-driven test from Step 6, but add these imports at the top:

```go
package usecases

import (
	"testing"

	"github.com/moneymate-2026/moneymate-backend/services/rewards/internal/domain"
)
```

Then run:

```powershell
go test .\services\rewards\internal\usecases
```

Expected result:

```text
ok github.com/moneymate-2026/moneymate-backend/services/rewards/internal/usecases
```

### Phase 15: Update Root Workspace Files

Edit file:

```text
go.work
```

Add inside `use (...)`:

```text
./services/rewards
```

Edit file:

```text
docker-compose.yml
```

Add the `rewards:` service block from Step 2.

In the existing `gateway:` service environment, add:

```yaml
- SERVICES_REWARDS_ADDR=rewards:9096
```

Edit file:

```text
docker-compose.prod.yml
```

Add a production rewards service similar to payment/notification:

```yaml
  rewards:
    image: ghcr.io/aamir-sufiyan-n/moneymate-rewards:latest
    env_file:
      - .env
    environment:
      - CONFIG_PATH=/app/config/config.yaml
      - POSTGRES_HOST=${POSTGRES_HOST}
      - POSTGRES_PORT=5432
      - REWARDS_DB_USER=${REWARDS_DB_USER}
      - REWARDS_DB_PASSWORD=${REWARDS_DB_PASSWORD}
      - ENVIRONMENT=${ENVIRONMENT}
      - JWT_ACCESS_SECRET=${JWT_ACCESS_SECRET}
      - JWT_REFRESH_SECRET=${JWT_REFRESH_SECRET}
    expose:
      - "9096"
    ports:
      - "9096:9096"
    restart: unless-stopped
    networks:
      - moneymate-network
```

In the production gateway environment, add:

```yaml
- SERVICES_REWARDS_ADDR=rewards:9096
```

Edit file:

```text
Taskfile.yml
```

Add this var under `vars:`:

```yaml
REWARDS_MIGRATIONS: /moneymate-backend/services/rewards/migrations
```

Add these tasks:

```yaml
  logs:rewards:
    desc: "Tail logs from rewards service only"
    cmd: docker compose -f {{.DEV_COMPOSE}} logs -f rewards

  migrate:rewards:
    desc: "Run rewards service migrations up"
    cmd: >
      docker compose -f {{.DEV_COMPOSE}} run --rm rewards
      migrate
      -path {{.REWARDS_MIGRATIONS}}
      -database "${REWARDS_DB_URL}"
      up

  migrate:rewards:down:
    desc: "Roll back last rewards migration"
    cmd: >
      docker compose -f {{.DEV_COMPOSE}} run --rm rewards
      migrate
      -path {{.REWARDS_MIGRATIONS}}
      -database "${REWARDS_DB_URL}"
      down 1

  migrate:rewards:status:
    desc: "Show current rewards migration version"
    cmd: >
      docker compose -f {{.DEV_COMPOSE}} run --rm rewards
      migrate
      -path {{.REWARDS_MIGRATIONS}}
      -database "${REWARDS_DB_URL}"
      version

  sqlc:rewards:
    desc: "Generate SQLC code for rewards service"
    cmd: >
      docker run --rm
      -v $(pwd)/services/rewards:/src
      -w /src
      sqlc/sqlc generate -f sqlc/sqlc.yml

  test:rewards:
    desc: "Run rewards service tests"
    dir: services/rewards
    cmd: go test ./...

  tidy:rewards:
    desc: "Run go mod tidy for rewards service"
    dir: services/rewards
    cmd: go mod tidy
```

Also update existing aggregate tasks:

```text
migrate:all should include migrate:rewards
test:all should include test:rewards
tidy:all should include tidy:rewards
```

Edit file:

```text
.env
```

Add:

```dotenv
REWARDS_DB_USER=rewards_user
REWARDS_DB_PASSWORD=rewards_password
REWARDS_DB_URL=postgres://rewards_user:rewards_password@postgres:5432/moneymate?sslmode=disable&search_path=rewards
```

### Phase 16: Update Gateway

Edit file:

```text
services/gateway/config/config.go
```

In `ServicesConfig`, add:

```go
RewardsAddr string `mapstructure:"rewards_addr"`
```

In `LoadConfig`, add:

```go
v.BindEnv("services.rewards_addr", "SERVICES_REWARDS_ADDR")
```

In the `cfg.Services.Downstream` map, add:

```go
"rewards": cfg.Services.RewardsAddr,
```

Edit file:

```text
services/gateway/config/config.yaml
```

Under `services:`, add:

```yaml
rewards_addr: "rewards:9096"
```

Create file:

```text
services/gateway/internal/transport/http/reward_routes.go
```

Paste the gateway reward route code from Step 11.

Edit file:

```text
services/gateway/internal/transport/http/router.go
```

Add this line near the other route registration calls:

```go
registerRewardRoutes(api, cfg.AuthMiddleware, cfg.Registry)
```

Edit file:

```text
services/gateway/internal/transport/http/admin_routes.go
```

Add a new function:

```go
func registerAdminRewardRuleRoutes(admin fiber.Router, registry *proxy.ServiceRegistry) {
	rules := admin.Group("/rewards/rules")
	rules.Post("/", proxy.HTTPProxy(registry, "rewards", "/admin/rewards/rules"))
	rules.Get("/", proxy.HTTPProxy(registry, "rewards", "/admin/rewards/rules"))
	rules.Get("/:id", proxy.HTTPProxy(registry, "rewards", "/admin/rewards/rules/:id"))
	rules.Put("/:id", proxy.HTTPProxy(registry, "rewards", "/admin/rewards/rules/:id"))
	rules.Patch("/:id/deactivate", proxy.HTTPProxy(registry, "rewards", "/admin/rewards/rules/:id/deactivate"))
}
```

Then call it inside `registerAdminRoutes`, after admin middleware is applied:

```go
registerAdminRewardRuleRoutes(admin, registry)
```

Do not remove the existing merchant admin rewards routes yet. Those are for merchant dashboard rewards history/summary and are separate.

### Phase 17: Boot Test The Skeleton

Run:

```powershell
docker compose up rewards --build
```

In another terminal, run:

```powershell
curl http://localhost:9096/health
curl http://localhost:9096/ready
```

Expected:

```json
{"service":"rewards","status":"ok"}
```

and:

```json
{"service":"rewards","status":"ready"}
```

If `/ready` fails on Kafka:

```text
Check KAFKA_BROKER, KAFKA_USERNAME, KAFKA_PASSWORD, and KAFKA_CA_CERT_PATH in .env.
The rewards config loader requires Kafka because this service consumes payment events.
```

If DB fails:

```text
Check POSTGRES_HOST, POSTGRES_DB, POSTGRES_SSL, REWARDS_DB_USER, and REWARDS_DB_PASSWORD.
Also make sure postgres is running and healthy.
```

### Phase 18: Then Add The Real Business Logic

After the skeleton boots:

```text
1. Generate SQLC.
2. Implement repo methods in services/rewards/internal/adapter/postgres/repo.
3. Implement usecases for rule CRUD and payout processing.
4. Replace placeholder routes with handlers.
5. Wire Kafka messages into payout processing.
6. Keep using fake PaymentClient until payment adds the real endpoint.
```

Suggested files to create for business logic:

```text
services/rewards/internal/adapter/postgres/repo/reward_repo.go
services/rewards/internal/adapter/paymentclient/fake.go
services/rewards/internal/usecases/rule_usecase.go
services/rewards/internal/usecases/payout_usecase.go
services/rewards/internal/transport/http/dto.go
services/rewards/internal/transport/http/rule_handler.go
services/rewards/internal/transport/http/payout_handler.go
services/rewards/internal/infra/kafkaconsumer/payment_completed.go
services/rewards/cmd/publish_fake_payment_event/main.go
```

Each file has one job:

```text
reward_repo.go
  Converts SQLC generated rows into domain structs.
  Handles duplicate insert errors for idempotency.

fake.go
  Fake PaymentClient for local dev and tests.

rule_usecase.go
  Validates admin inputs and calls repository rule methods.

payout_usecase.go
  Loads active rule, calculates reward, inserts payout row, calls PaymentClient, marks completed/failed.

dto.go
  Request/response structs for HTTP JSON.

rule_handler.go
  Admin CRUD HTTP handlers for reward rules.

payout_handler.go
  GET /rewards/me and GET /rewards?transaction_id=...

payment_completed.go
  Kafka event struct and JSON parsing.

publish_fake_payment_event/main.go
  Local test publisher until payment emits the real event.
```

## Step 1: Service Skeleton

Create this folder structure:

```text
services/rewards/
  .air.toml
  Dockerfile
  go.mod
  cmd/main.go
  config/config.go
  config/config.yaml
  internal/app/app.go
  internal/adapter/postgres/db.go
  internal/adapter/postgres/repo/
  internal/adapter/paymentclient/
  internal/domain/
  internal/infra/kafkaconsumer/
  internal/transport/http/
  internal/usecases/
  migrations/
  sqlc/sqlc.yml
  sqlc/queries/
  sqlc/generated/
```

Recommended module:

```go
module github.com/moneymate-2026/moneymate-backend/services/rewards

go 1.26.4

require (
	github.com/gofiber/fiber/v3 v3.4.0
	github.com/golang-migrate/migrate/v4 v4.19.1
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.10.0
	github.com/joho/godotenv v1.5.1
	github.com/moneymate-2026/moneymate-backend/shared v0.0.0
	github.com/spf13/viper v1.21.0
)
```

Use this service config:

```yaml
server:
  http_addr: "9096"

database:
  max_open_conns: 25
  min_open_conns: 5
  max_conn_lifetime: 15m
  max_idle_time: 5m
  migrations_path: "/app/migrations"

jwt:
  access_expiry_minutes: 15
  refresh_expiry_hours: 720
  tx_token_expiry_secs: 60

rewards:
  payment_completed_topic: "moneymate.payment.completed"
  consumer_group: "moneymate-rewards-svc"
  fake_payment_client: true
```

Config loader:

```go
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	sharedconfig "github.com/moneymate-2026/moneymate-backend/shared/config"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	HTTPAddr string `mapstructure:"http_addr"`
}

type RewardsConfig struct {
	PaymentCompletedTopic string `mapstructure:"payment_completed_topic"`
	ConsumerGroup         string `mapstructure:"consumer_group"`
	FakePaymentClient     bool   `mapstructure:"fake_payment_client"`
}

type Config struct {
	Env                   string
	Server                ServerConfig `mapstructure:"server"`
	Rewards              RewardsConfig `mapstructure:"rewards"`
	Database              sharedconfig.DatabaseConfig
	Kafka                 sharedconfig.KafkaConfig
	JWT                   sharedconfig.JWTConfig
	InternalServiceSecret string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	yamlPath := os.Getenv("CONFIG_PATH")
	if yamlPath == "" {
		yamlPath = "./config/config.yaml"
	}

	v := viper.New()
	v.SetConfigFile(yamlPath)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("read config: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	cfg.Database = sharedconfig.LoadDatabaseConfig(v, "rewards")
	cfg.Kafka = sharedconfig.LoadKafkaConfig(v)
	cfg.JWT = sharedconfig.LoadJWTConfig(v)
	cfg.Env = sharedconfig.Get("ENVIRONMENT", "dev")
	cfg.InternalServiceSecret = sharedconfig.MustGet("INTERNAL_SERVICE_SECRET")

	return &cfg, nil
}
```

Use the payment or notification `internal/adapter/postgres/db.go` as the template, but import the rewards config package and set `search_path=rewards`.

Add `/health` and `/ready` in `internal/app/app.go`.

Minimum readiness:

```go
server.Get("/health", func(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok", "service": "rewards"})
})

server.Get("/ready", func(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status": "unavailable",
			"dependency": "postgres",
			"error": err.Error(),
		})
	}

	if kafkaConsumer == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status": "unavailable",
			"dependency": "kafka",
			"error": "consumer not initialized",
		})
	}

	return c.JSON(fiber.Map{"status": "ready", "service": "rewards"})
})
```

## Step 2: Workspace, Docker, Taskfile, Env

Add rewards to `go.work`:

```text
./services/rewards
```

Important: every existing service Dockerfile copies all workspace module `go.mod` files before `go work sync`. After adding rewards to `go.work`, update the base stage in all service Dockerfiles to include:

```dockerfile
COPY services/rewards/go.mod services/rewards/go.mod
```

Add a rewards service to `docker-compose.yml`:

```yaml
  rewards:
    build:
      context: .
      dockerfile: ./services/rewards/Dockerfile
      target: dev
    volumes:
      - ./:/moneymate-backend
      - ./shared:/moneymate-backend/shared
      - ./services/rewards/migrations:/app/migrations
    working_dir: /moneymate-backend/services/rewards
    env_file:
      - .env
    environment:
      - CONFIG_PATH=/moneymate-backend/services/rewards/config/config.yaml
      - POSTGRES_HOST=${POSTGRES_HOST}
      - POSTGRES_PORT=5432
      - ENVIRONMENT=dev
      - JWT_ACCESS_SECRET=${JWT_ACCESS_SECRET}
      - JWT_REFRESH_SECRET=${JWT_REFRESH_SECRET}
      - REWARDS_DB_USER=${REWARDS_DB_USER}
      - REWARDS_DB_PASSWORD=${REWARDS_DB_PASSWORD}
    ports:
      - "9096:9096"
    restart: unless-stopped
    depends_on:
      postgres:
        condition: service_healthy
    networks:
      - moneymate-network
```

Add this to gateway environment:

```yaml
- SERVICES_REWARDS_ADDR=rewards:9096
```

Add these `.env` values:

```dotenv
REWARDS_DB_USER=rewards_user
REWARDS_DB_PASSWORD=rewards_password
REWARDS_DB_URL=postgres://rewards_user:rewards_password@postgres:5432/moneymate?sslmode=disable&search_path=rewards
```

Add Taskfile vars and tasks:

```yaml
vars:
  REWARDS_MIGRATIONS: /moneymate-backend/services/rewards/migrations

tasks:
  logs:rewards:
    desc: "Tail logs from rewards service only"
    cmd: docker compose -f {{.DEV_COMPOSE}} logs -f rewards

  migrate:rewards:
    desc: "Run rewards service migrations up"
    cmd: >
      docker compose -f {{.DEV_COMPOSE}} run --rm rewards
      migrate
      -path {{.REWARDS_MIGRATIONS}}
      -database "${REWARDS_DB_URL}"
      up

  sqlc:rewards:
    desc: "Generate SQLC code for rewards service"
    cmd: >
      docker run --rm
      -v $(pwd)/services/rewards:/src
      -w /src
      sqlc/sqlc generate -f sqlc/sqlc.yml

  test:rewards:
    desc: "Run rewards service tests"
    dir: services/rewards
    cmd: go test ./...
```

Also add `migrate:rewards`, `sqlc:rewards`, `test:rewards`, and `tidy:rewards` to the existing aggregate tasks where appropriate.

## Step 3: Schema And Migrations

Create:

```text
services/rewards/migrations/000001_create_rewards_tables.up.sql
services/rewards/migrations/000001_create_rewards_tables.down.sql
```

Recommended DDL:

```sql
CREATE TYPE reward_recipient_type AS ENUM ('user', 'merchant');
CREATE TYPE reward_payout_status AS ENUM ('pending', 'completed', 'failed');

CREATE TABLE reward_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    min_percentage_bps INT NOT NULL CHECK (min_percentage_bps >= 0),
    max_percentage_bps INT NOT NULL CHECK (max_percentage_bps >= min_percentage_bps),
    min_transaction_amount_paise BIGINT NOT NULL DEFAULT 0 CHECK (min_transaction_amount_paise >= 0),
    max_payout_amount_paise BIGINT NOT NULL CHECK (max_payout_amount_paise > 0),
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_reward_rules_one_active
    ON reward_rules (active)
    WHERE active = TRUE;

CREATE TABLE reward_payouts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID NOT NULL,
    recipient_id UUID NOT NULL,
    recipient_account_id UUID NOT NULL,
    recipient_type reward_recipient_type NOT NULL,
    rule_id UUID REFERENCES reward_rules(id),
    original_amount_paise BIGINT NOT NULL CHECK (original_amount_paise > 0),
    reward_percentage_bps INT NOT NULL CHECK (reward_percentage_bps >= 0),
    reward_amount_paise BIGINT NOT NULL CHECK (reward_amount_paise >= 0),
    status reward_payout_status NOT NULL DEFAULT 'pending',
    payment_transaction_id UUID,
    failure_reason TEXT,
    event_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    UNIQUE (transaction_id, recipient_type)
);

CREATE INDEX idx_reward_payouts_recipient_created
    ON reward_payouts (recipient_id, created_at DESC);

CREATE INDEX idx_reward_payouts_transaction_id
    ON reward_payouts (transaction_id);

CREATE INDEX idx_reward_payouts_status
    ON reward_payouts (status);
```

Down migration:

```sql
DROP INDEX IF EXISTS idx_reward_payouts_status;
DROP INDEX IF EXISTS idx_reward_payouts_transaction_id;
DROP INDEX IF EXISTS idx_reward_payouts_recipient_created;
DROP TABLE IF EXISTS reward_payouts;

DROP INDEX IF EXISTS idx_reward_rules_one_active;
DROP TABLE IF EXISTS reward_rules;

DROP TYPE IF EXISTS reward_payout_status;
DROP TYPE IF EXISTS reward_recipient_type;
```

Why basis points:

```text
1 basis point = 0.01 percent
150 bps = 1.50 percent
reward = amount_paise * bps / 10000
```

This keeps reward math deterministic and avoids float rounding bugs.

## Step 4: SQLC

Create `services/rewards/sqlc/sqlc.yml`:

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "queries"
    schema: "../migrations"
    gen:
      go:
        package: "generated"
        out: "generated"
        sql_package: "pgx/v5"
        emit_interface: true
        emit_empty_slices: true
        emit_json_tags: false
        emit_pointers_for_null_types: true
        overrides:
          - db_type: uuid
            go_type: github.com/google/uuid.UUID
          - db_type: timestamptz
            go_type: time.Time
```

Create `services/rewards/sqlc/queries/reward_rules.sql`:

```sql
-- name: CreateRewardRule :one
INSERT INTO reward_rules (
    name, min_percentage_bps, max_percentage_bps,
    min_transaction_amount_paise, max_payout_amount_paise,
    active, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: ListRewardRules :many
SELECT * FROM reward_rules
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetRewardRuleByID :one
SELECT * FROM reward_rules
WHERE id = $1;

-- name: GetActiveRewardRule :one
SELECT * FROM reward_rules
WHERE active = TRUE
ORDER BY created_at DESC
LIMIT 1;

-- name: UpdateRewardRule :one
UPDATE reward_rules
SET name = $2,
    min_percentage_bps = $3,
    max_percentage_bps = $4,
    min_transaction_amount_paise = $5,
    max_payout_amount_paise = $6,
    active = $7,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeactivateRewardRule :one
UPDATE reward_rules
SET active = FALSE,
    updated_at = NOW()
WHERE id = $1
RETURNING *;
```

Create `services/rewards/sqlc/queries/reward_payouts.sql`:

```sql
-- name: InsertRewardPayout :one
INSERT INTO reward_payouts (
    transaction_id, recipient_id, recipient_account_id, recipient_type,
    rule_id, original_amount_paise, reward_percentage_bps,
    reward_amount_paise, status, event_payload
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, 'pending', $9
) RETURNING *;

-- name: GetRewardPayoutByID :one
SELECT * FROM reward_payouts
WHERE id = $1;

-- name: GetRewardPayoutByOriginalTransaction :many
SELECT * FROM reward_payouts
WHERE transaction_id = $1
ORDER BY created_at DESC;

-- name: ListRewardPayoutsByRecipient :many
SELECT * FROM reward_payouts
WHERE recipient_id = $1
  AND (sqlc.narg('status')::reward_payout_status IS NULL OR status = sqlc.narg('status')::reward_payout_status)
ORDER BY created_at DESC
LIMIT sqlc.arg('limit_count') OFFSET sqlc.arg('offset_count');

-- name: MarkRewardPayoutCompleted :one
UPDATE reward_payouts
SET status = 'completed',
    payment_transaction_id = $2,
    completed_at = NOW(),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: MarkRewardPayoutFailed :one
UPDATE reward_payouts
SET status = 'failed',
    failure_reason = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;
```

Run:

```bash
task sqlc:rewards
```

## Step 5: Domain Contracts

Create `internal/domain/reward.go`:

```go
package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RecipientType string

const (
	RecipientTypeUser     RecipientType = "user"
	RecipientTypeMerchant RecipientType = "merchant"
)

type PayoutStatus string

const (
	PayoutStatusPending   PayoutStatus = "pending"
	PayoutStatusCompleted PayoutStatus = "completed"
	PayoutStatusFailed    PayoutStatus = "failed"
)

type RewardRule struct {
	ID                         uuid.UUID
	Name                       string
	MinPercentageBPS           int32
	MaxPercentageBPS           int32
	MinTransactionAmountPaise  int64
	MaxPayoutAmountPaise       int64
	Active                     bool
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
}

type RewardPayout struct {
	ID                  uuid.UUID
	TransactionID       uuid.UUID
	RecipientID         uuid.UUID
	RecipientAccountID  uuid.UUID
	RecipientType       RecipientType
	RuleID              *uuid.UUID
	OriginalAmountPaise int64
	RewardPercentageBPS int32
	RewardAmountPaise   int64
	Status              PayoutStatus
	PaymentTransactionID *uuid.UUID
	FailureReason       *string
	CreatedAt           time.Time
	CompletedAt         *time.Time
}

type RewardRepository interface {
	CreateRule(ctx context.Context, rule RewardRule, createdBy *uuid.UUID) (*RewardRule, error)
	ListRules(ctx context.Context, limit, offset int32) ([]*RewardRule, error)
	GetRule(ctx context.Context, id uuid.UUID) (*RewardRule, error)
	GetActiveRule(ctx context.Context) (*RewardRule, error)
	UpdateRule(ctx context.Context, rule RewardRule) (*RewardRule, error)
	DeactivateRule(ctx context.Context, id uuid.UUID) (*RewardRule, error)

	InsertPayout(ctx context.Context, payout RewardPayout, eventPayload []byte) (*RewardPayout, error)
	ListPayoutsByRecipient(ctx context.Context, recipientID uuid.UUID, status *PayoutStatus, limit, offset int32) ([]*RewardPayout, error)
	ListPayoutsByTransaction(ctx context.Context, transactionID uuid.UUID) ([]*RewardPayout, error)
	MarkCompleted(ctx context.Context, payoutID, paymentTransactionID uuid.UUID) (*RewardPayout, error)
	MarkFailed(ctx context.Context, payoutID uuid.UUID, reason string) (*RewardPayout, error)
}
```

Create `internal/domain/payment.go`:

```go
package domain

import (
	"context"

	"github.com/google/uuid"
)

type PaymentClient interface {
	ExecuteRewardPayout(ctx context.Context, recipientAccountID uuid.UUID, amountPaise int64) (txID uuid.UUID, err error)
}
```

## Step 6: Pure Rule Engine

Create `internal/usecases/rule_engine.go`:

```go
package usecases

import "github.com/moneymate-2026/moneymate-backend/services/rewards/internal/domain"

type RollBPSFunc func(minBPS, maxBPS int32) int32

type RewardCalculation struct {
	Eligible          bool
	PercentageBPS     int32
	RewardAmountPaise int64
	Reason            string
}

func CalculateReward(amountPaise int64, rule domain.RewardRule, roll RollBPSFunc) RewardCalculation {
	if !rule.Active {
		return RewardCalculation{Eligible: false, Reason: "rule inactive"}
	}
	if amountPaise < rule.MinTransactionAmountPaise {
		return RewardCalculation{Eligible: false, Reason: "below threshold"}
	}
	if amountPaise <= 0 {
		return RewardCalculation{Eligible: false, Reason: "invalid amount"}
	}

	bps := roll(rule.MinPercentageBPS, rule.MaxPercentageBPS)
	if bps < rule.MinPercentageBPS {
		bps = rule.MinPercentageBPS
	}
	if bps > rule.MaxPercentageBPS {
		bps = rule.MaxPercentageBPS
	}

	reward := amountPaise * int64(bps) / 10000
	if reward > rule.MaxPayoutAmountPaise {
		reward = rule.MaxPayoutAmountPaise
	}
	if reward <= 0 {
		return RewardCalculation{Eligible: false, PercentageBPS: bps, Reason: "zero reward"}
	}

	return RewardCalculation{
		Eligible:          true,
		PercentageBPS:     bps,
		RewardAmountPaise: reward,
	}
}
```

Table-driven tests:

```go
func TestCalculateReward(t *testing.T) {
	base := domain.RewardRule{
		MinPercentageBPS:          0,
		MaxPercentageBPS:          150,
		MinTransactionAmountPaise: 10000,
		MaxPayoutAmountPaise:      5000,
		Active:                    true,
	}

	tests := []struct {
		name       string
		amount     int64
		rule       domain.RewardRule
		rolledBPS  int32
		eligible   bool
		wantAmount int64
	}{
		{name: "below threshold", amount: 9999, rule: base, rolledBPS: 100, eligible: false},
		{name: "normal range", amount: 100000, rule: base, rolledBPS: 100, eligible: true, wantAmount: 1000},
		{name: "at cap", amount: 1000000, rule: base, rolledBPS: 150, eligible: true, wantAmount: 5000},
		{name: "inactive", amount: 100000, rule: func() domain.RewardRule { r := base; r.Active = false; return r }(), rolledBPS: 100, eligible: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateReward(tt.amount, tt.rule, func(_, _ int32) int32 { return tt.rolledBPS })
			if got.Eligible != tt.eligible {
				t.Fatalf("eligible = %v, want %v", got.Eligible, tt.eligible)
			}
			if got.RewardAmountPaise != tt.wantAmount {
				t.Fatalf("amount = %d, want %d", got.RewardAmountPaise, tt.wantAmount)
			}
		})
	}
}
```

## Step 7: Admin CRUD For Reward Rules

HTTP request DTO:

```go
type RewardRuleRequest struct {
	Name                      string `json:"name"`
	MinPercentageBPS          int32  `json:"min_percentage_bps"`
	MaxPercentageBPS          int32  `json:"max_percentage_bps"`
	MinTransactionAmountPaise int64  `json:"min_transaction_amount_paise"`
	MaxPayoutAmountPaise      int64  `json:"max_payout_amount_paise"`
	Active                    bool   `json:"active"`
}
```

Validation rules:

```text
name is required
min_percentage_bps >= 0
max_percentage_bps >= min_percentage_bps
max_percentage_bps <= 10000
min_transaction_amount_paise >= 0
max_payout_amount_paise > 0
only one active rule is allowed
```

Auth recommendation:

```text
Gateway:
  /api/v1/admin/... already applies authMiddleware + RequireRole("admin").

Rewards service:
  Either parse JWT again with shared/pkg/jwt for direct-service safety,
  or accept that admin routes are only exposed through the gateway in dev.

Internal routes:
  Use payment/auth style X-Internal-Secret middleware.
```

Do not copy merchant's no-op admin auth into this service.

## Step 8: Kafka Consumer Skeleton

Use a provisional event shape until payment confirms the final one:

```go
type PaymentCompletedEvent struct {
	EventID            string    `json:"event_id"`
	EventType          string    `json:"event_type"`
	TransactionID      uuid.UUID `json:"transaction_id"`
	RecipientID        uuid.UUID `json:"recipient_id"`
	RecipientAccountID uuid.UUID `json:"recipient_account_id"`
	RecipientType      string    `json:"recipient_type"`
	AmountPaise        int64     `json:"amount_paise"`
	OccurredAt         time.Time `json:"occurred_at"`
}
```

Consumer flow:

```text
1. Unmarshal event payload.
2. Load active reward rule.
3. Calculate reward using pure rule engine.
4. If not eligible, return nil so Kafka can move on.
5. Insert reward_payouts with status=pending.
6. If unique violation on transaction_id + recipient_type, treat as already processed and return nil.
7. Call PaymentClient.ExecuteRewardPayout.
8. Mark payout completed with payment tx id.
9. If payout fails, mark reward_payouts failed with reason.
```

Idempotency handling:

```go
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
```

The unique constraint is:

```sql
UNIQUE (transaction_id, recipient_type)
```

That protects against at-least-once Kafka delivery.

## Step 9: PaymentClient Interface

Create a fake client first:

```go
package paymentclient

import (
	"context"
	"log"

	"github.com/google/uuid"
)

type FakeClient struct{}

func NewFakeClient() *FakeClient {
	return &FakeClient{}
}

func (c *FakeClient) ExecuteRewardPayout(ctx context.Context, recipientAccountID uuid.UUID, amountPaise int64) (uuid.UUID, error) {
	txID := uuid.New()
	log.Printf("fake reward payout: recipient_account_id=%s amount_paise=%d tx_id=%s", recipientAccountID, amountPaise, txID)
	return txID, nil
}
```

Later, when payment adds the internal endpoint, replace only this adapter:

```go
type HTTPClient struct {
	BaseURL        string
	InternalSecret string
	HTTP           *http.Client
}

func (c *HTTPClient) ExecuteRewardPayout(ctx context.Context, recipientAccountID uuid.UUID, amountPaise int64) (uuid.UUID, error) {
	// POST /internal/payment/reward-payouts
	// Headers:
	//   X-Internal-Secret: ...
	// Body:
	//   recipient_account_id
	//   amount_paise
	// Return:
	//   payment_transaction_id
	return uuid.Nil, nil
}
```

Nothing in the usecase or consumer should change when this swap happens.

## Step 10: Query Endpoints

`GET /rewards/me`

```text
Headers:
  Authorization: Bearer <access token>

Query:
  limit=50
  offset=0
  status=pending|completed|failed

Response:
  {
    "success": true,
    "data": [
      {
        "id": "...",
        "transaction_id": "...",
        "reward_amount_paise": 125,
        "reward_amount": "1.25",
        "status": "completed",
        "created_at": "..."
      }
    ]
  }
```

`GET /rewards?transaction_id=<uuid>`

```text
Headers:
  Authorization: Bearer <access token>

Response:
  Rewards generated from the original payment transaction.
```

Use `shared/pkg/money.FormatPaise` for display strings.

## Step 11: Gateway Wiring

Update gateway config types:

```go
type ServicesConfig struct {
	AuthAddr     string            `mapstructure:"auth_addr"`
	MerchantAddr string            `mapstructure:"merch_addr"`
	PaymentAddr  string            `mapstructure:"payment_addr"`
	SupportAddr  string            `mapstructure:"support_addr"`
	RewardsAddr  string            `mapstructure:"rewards_addr"`
	Downstream   map[string]string `mapstructure:"downstream"`
}
```

Bind env:

```go
v.BindEnv("services.rewards_addr", "SERVICES_REWARDS_ADDR")
```

Add downstream registry entry:

```go
cfg.Services.Downstream = map[string]string{
	"payment": cfg.Services.PaymentAddr,
	"support": cfg.Services.SupportAddr,
	"rewards": cfg.Services.RewardsAddr,
}
```

Add `services/gateway/internal/transport/http/reward_routes.go`:

```go
package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/moneymate-2026/moneymate-backend/gateway/internal/proxy"
)

func registerRewardRoutes(api fiber.Router, authMiddleware fiber.Handler, registry *proxy.ServiceRegistry) {
	rewards := api.Group("/rewards")
	rewards.Use(authMiddleware)
	rewards.Get("/me", proxy.HTTPProxy(registry, "rewards", "/rewards/me"))
	rewards.Get("/", proxy.HTTPProxy(registry, "rewards", "/rewards"))
}
```

Call it from `router.go`:

```go
registerRewardRoutes(api, cfg.AuthMiddleware, cfg.Registry)
```

Add admin rule proxying inside `admin_routes.go` after admin middleware is applied:

```go
registerAdminRewardRuleRoutes(admin, registry)
```

## Step 12: Local Fake Publisher

Until payment emits the final event, add:

```text
services/rewards/cmd/publish_fake_payment_event/main.go
```

Flow:

```text
1. Load rewards config.
2. Create shared/pkg/kafka Producer.
3. Marshal PaymentCompletedEvent.
4. Publish to cfg.Rewards.PaymentCompletedTopic.
```

Example payload:

```json
{
  "event_id": "local-test-001",
  "event_type": "moneymate.payment.completed",
  "transaction_id": "00000000-0000-0000-0000-000000000001",
  "recipient_id": "00000000-0000-0000-0000-000000000002",
  "recipient_account_id": "00000000-0000-0000-0000-000000000003",
  "recipient_type": "user",
  "amount_paise": 250000,
  "occurred_at": "2026-08-20T00:00:00Z"
}
```

## Acceptance Checklist

Task 1 complete when:

```text
docker compose up rewards
curl http://localhost:9096/health
curl http://localhost:9096/ready
```

returns healthy JSON and startup logs show migrations ran.

Task 2 complete when:

```text
task migrate:rewards
task sqlc:rewards
go test ./services/rewards/...
```

all pass.

Task 3 complete when:

```text
rule engine tests cover:
  below threshold
  normal reward
  capped reward
  inactive rule
  invalid amount
```

Task 4 complete when:

```text
admin can create/list/update/deactivate rules through gateway
only one active rule can exist
invalid percentage/cap inputs return 400
```

Task 5 complete when:

```text
fake publisher emits an event
consumer receives it
reward_payouts row is inserted as pending
publishing the same event twice does not create a duplicate
```

Task 6 complete when:

```text
fake PaymentClient is called
pending payout becomes completed
fake client failure path marks payout failed
```

Task 7 complete when:

```text
GET /api/v1/rewards/me returns current user's payouts
GET /api/v1/rewards?transaction_id=... returns transaction-linked payouts
empty state returns []
```

Task 8 complete when:

```text
gateway routes reach rewards service
docker compose starts gateway + rewards together
Taskfile has migrate/sqlc/test/logs/tidy rewards tasks
```

## Known Blockers

Only one part is blocked:

```text
Real money movement for ExecuteRewardPayout.
```

Payment needs to add an internal endpoint that can credit a reward amount into a user or merchant account and return a payment transaction id.

Until then:

```text
Keep PaymentClient as an interface.
Use FakeClient locally and in tests.
Persist reward_payouts exactly as if the payout happened.
Swap only the adapter later.
```

## Design Decisions To Keep Stable

Use `services/rewards`, not `services/reward`, because the database bootstrap and gateway generic route list already use `rewards`.

Use paise everywhere in code. Convert to display strings only in HTTP responses.

Keep the rule engine pure. No DB, Kafka, HTTP, randomness source, or time reads inside the calculation function.

Do not reuse merchant reward tables. Merchant rewards are store dashboard and redemption data; this service calculates cashback from payment events.

Do not depend on payment tables. The service owns `reward_rules` and `reward_payouts`; payment communication happens only through Kafka input and `PaymentClient` output.

Treat Kafka as at-least-once. The unique constraint plus duplicate handling is mandatory.

## File Edit Map

When implementation begins, expect to create or edit these files:

```text
Create:
  services/rewards/**

Edit:
  go.work
  docker-compose.yml
  docker-compose.prod.yml
  Taskfile.yml
  services/gateway/config/config.go
  services/gateway/config/config.yaml
  services/gateway/cmd/main.go only if RouteConfig needs new fields
  services/gateway/internal/transport/http/router.go
  services/gateway/internal/transport/http/admin_routes.go
  services/gateway/internal/transport/http/reward_routes.go
  services/*/Dockerfile base stages to copy services/rewards/go.mod

Probably no change needed:
  infra/postgres/bootstrap/*.sql because rewards schema and role already exist
  shared/pkg/money because paise helpers already exist
  shared/pkg/kafka because producer/consumer wrappers already exist
```

## Suggested First Commit Boundary

Commit 1:

```text
services/rewards skeleton
config
Dockerfile
go.work
docker-compose
Taskfile
/health
/ready
```

Commit 2:

```text
migrations
sqlc queries
repo mappings
rule engine
rule engine tests
```

Commit 3:

```text
admin rule CRUD
query endpoints
gateway routes
```

Commit 4:

```text
Kafka consumer
fake publisher
PaymentClient interface
fake payout client
idempotency tests
```
