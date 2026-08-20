package http

import (
	"github.com/gofiber/fiber/v3"

	"github.com/moneymate-2026/moneymate-backend/services/rewards/internal/usecases"
	sharedjwt "github.com/moneymate-2026/moneymate-backend/shared/pkg/jwt"
)

type RewardHandler struct {
	rewardUC usecases.RewardUsecase
}

func NewRewardHandler(rewardUC usecases.RewardUsecase) *RewardHandler {
	return &RewardHandler{rewardUC: rewardUC}
}

func RegisterRoutes(router fiber.Router, rh *RewardHandler, internalSecret string) {
	jwtCfg := sharedjwt.Config{
		AccessSecret:     resolveEnv("JWT_ACCESS_SECRET"),
		RefreshSecret:    resolveEnv("JWT_REFRESH_SECRET"),
		AccessExpiryMins: 15,
		RefreshExpiryHrs: 720,
	}

	api := router.Group("/rewards")
	api.Use(RequireUserID(jwtCfg))

	api.Get("/me", rh.ListMyPayouts)
	api.Get("/", rh.ListPayoutsByTransaction)

	admin := router.Group("/admin/rewards/rules")
	admin.Use(RequireInternalSecret(internalSecret))

	admin.Post("/", rh.CreateRule)
	admin.Get("/", rh.ListRules)
	admin.Get("/:id", rh.GetRule)
	admin.Put("/:id", rh.UpdateRule)
	admin.Patch("/:id/deactivate", rh.DeactivateRule)

	internal := router.Group("/internal", RequireInternalSecret(internalSecret))
	internal.Get("/rewards/payouts/:id", rh.GetPayoutByID)
	internal.Post("/rewards/replay-failed", rh.ReplayFailed)

	dev := router.Group("/dev", RequireInternalSecret(internalSecret))
	dev.Post("/fake-payment-event", rh.FakePaymentEvent)
}
