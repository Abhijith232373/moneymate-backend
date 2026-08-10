package http

import (
	"github.com/gofiber/fiber/v3"

	authclient "github.com/moneymate-2026/moneymate-backend/services/payment/internal/adapter/authClient"
	sharedjwt "github.com/moneymate-2026/moneymate-backend/shared/pkg/jwt"
)

func RegisterRoutes(router fiber.Router, wh *WalletHandler, th *TransferHandler, jwtCfg sharedjwt.Config, authClient *authclient.Client) {
	pay := router.Group("/payment", RequireUserID(jwtCfg))

	pay.Post("/wallets", wh.CreateWallet)
	pay.Get("/wallets/me", wh.GetMyWallet)
	pay.Get("/wallets/:id", wh.GetWalletByID)

	pay.Post("/transfers", RequireTransactionToken(authClient), th.Transfer)
	pay.Get("/transactions/:id", th.GetTransaction)
}