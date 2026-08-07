package http

import (
	"github.com/gofiber/fiber/v3"

	sharedjwt "github.com/moneymate-2026/moneymate-backend/shared/pkg/jwt"
)

// RegisterRoutes wires the /payment/... API. The gateway proxies
// /api/v1/payment/* here (that is the S1 team's job, done in parallel).
func RegisterRoutes(router fiber.Router, wh *WalletHandler, th *TransferHandler, jwtCfg sharedjwt.Config) {
	pay := router.Group("/payment", RequireUserID(jwtCfg))

	pay.Post("/wallets", wh.CreateWallet)
	pay.Get("/wallets/me", wh.GetMyWallet)
	pay.Get("/wallets/:id", wh.GetWalletByID)

	pay.Post("/transfers", RequireTransactionToken(jwtCfg), th.Transfer)
	pay.Get("/transactions/:id", th.GetTransaction)
}
