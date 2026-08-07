package http

import (
	"github.com/gofiber/fiber/v3"

	"github.com/moneymate-2026/moneymate-backend/auth/internal/domain"
)

type Handlers struct {
	Auth    *AuthHandler
	Role    *RoleHandler
	User    *UserHandler
	UserPin *UserPinHandler
}

func RegisterRoutes(router fiber.Router, h *Handlers, ) {
	registerAuthRoutes(router, h.Auth, )
	registerRoleRoutes(router, h.Role, )
	registerUserRoutes(router, h.User, )
	registerUserPinRoutes(router, h.UserPin, )
}

func registerAuthRoutes(router fiber.Router, h *AuthHandler) {
	auth := router.Group("/auth")
	auth.Get("/health", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ok",
			"service": "auth",
		})
	})
	auth.Post("/login", h.Login)
	auth.Post("/logout", RequireUserID, h.Logout)
	auth.Post("/otp/send", h.SendRegistrationOTP)
	auth.Post("/otp/verify", h.VerifyRegistrationOTP)
	auth.Post("/user/register", h.Register(domain.AccountTypeUser))
	auth.Post("/merchant/register", h.Register(domain.AccountTypeMerchant))

	internal := router.Group("/internal")
	internal.Post("/auth/verify-access-token", h.VerifyAccessToken)
	internal.Post("/auth/verify-transaction-token", h.VerifyTransactionToken)
	// internal.Post("/auth/verify-pin", h.VerifyPINInternal)
	internal.Get("/auth/users/:id", h.GetUserByID)
}

func registerRoleRoutes(router fiber.Router, h *RoleHandler ) {
	roles := router.Group("/admin/roles", )
	roles.Post("/", h.CreateRole)
	roles.Get("/", h.ListRoles)
	roles.Get("/:id", h.GetRole)
	roles.Put("/:id", h.UpdateRole)
	roles.Delete("/:id", h.DeleteRole)
	roles.Post("/assign", h.AssignRoleToUser)
	roles.Delete("/users/:userId/roles/:roleId", h.RemoveRoleFromUser)
	roles.Get("/users/:userId", h.GetUserRoles)
} 
func registerUserRoutes(router fiber.Router, h *UserHandler ) {
	users := router.Group("/admin/users", )
	users.Post("/",h.CreateUser)
	users.Get("/", h.ListUsers)
	users.Get("/:id", h.GetUser)
	users.Put("/:id", h.UpdateUser)
	users.Patch("/:id/status", h.UpdateUserStatus)
	users.Delete("/:id", h.DeleteUser)
}

func registerUserPinRoutes(router fiber.Router, h *UserPinHandler ) {
	pins := router.Group("/user/pin", RequireUserID)
	pins.Post("/", h.SetPIN)
	pins.Put("/", h.UpdatePIN)
	pins.Post("/verify", h.VerifyPIN)
}