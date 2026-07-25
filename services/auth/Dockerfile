
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

RUN go work sync

FROM base AS builder
WORKDIR /moneymate-backend
COPY shared/ shared/
COPY services/auth/ services/auth/
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/auth ./services/auth/cmd/main.go

FROM alpine:3.19 AS prod
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app

RUN mkdir -p /app/config
COPY --from=builder /bin/auth ./auth
COPY --from=builder /moneymate-backend/services/auth/config/config.yaml ./config/config.yaml
COPY --from=builder /moneymate-backend/services/auth/migrations ./migrations

USER appuser
EXPOSE 8081
CMD ["./auth"]