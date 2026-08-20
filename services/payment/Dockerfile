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
WORKDIR /moneymate-backend/services/payment
RUN go install github.com/air-verse/air@latest
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
CMD ["air", "-c", ".air.toml"]

FROM base AS builder
WORKDIR /moneymate-backend
COPY shared/ shared/
COPY services/payment/ services/payment/
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/payment ./services/payment/cmd/main.go

FROM alpine:3.19 AS prod
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app

RUN mkdir -p /app/config
COPY --from=builder /bin/payment ./payment
COPY --from=builder /moneymate-backend/services/payment/config/config.yaml ./config/config.yaml
COPY --from=builder /moneymate-backend/services/payment/migrations ./migrations
COPY --from=builder /moneymate-backend/shared/certs/aiven-ca.pem ./certs/aiven-ca.pem

USER appuser
EXPOSE 9094
CMD ["./payment"]