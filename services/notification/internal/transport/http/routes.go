package http

import (
	"github.com/gofiber/fiber/v3"

	"github.com/moneymate-2026/moneymate-backend/services/notification/internal/usecases"
	sharedjwt "github.com/moneymate-2026/moneymate-backend/shared/pkg/jwt"
)

func RegisterRoutes(router fiber.Router, du *usecases.DeviceUsecase, iu *usecases.InboxUsecase, pu *usecases.PreferenceUsecase, jwtCfg sharedjwt.Config) {
	deviceHandler := NewDeviceHandler(du)
	inboxHandler := NewInboxHandler(iu)
	preferenceHandler := NewPreferenceHandler(pu)

	n := router.Group("/notification", RequireRecipient(jwtCfg))

	n.Post("/devices", deviceHandler.Register)
	n.Delete("/devices/:device_id", deviceHandler.Revoke)

	n.Get("/inbox", inboxHandler.List)
	n.Get("/inbox/:id", inboxHandler.Get)
	n.Patch("/inbox/:id/read", inboxHandler.MarkRead)

	n.Get("/preferences", preferenceHandler.List)
	n.Put("/preferences", preferenceHandler.Upsert)
}
