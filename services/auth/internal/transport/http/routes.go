package http

import (
	"github.com/gofiber/fiber/v3"

	"github.com/moneymate-2026/moneymate-backend/auth/internal/domain"
)

type Handlers struct {
	Auth  *AuthHandler
	Role  *RoleHandler
	User  *UserHandler
}

func RegisterRoutes(router fiber.Router, h *Handlers, ) {
	registerAuthRoutes(router, h.Auth, )
	registerRoleRoutes(router, h.Role, )
	registerUserRoutes(router, h.User, )
}

func registerAuthRoutes(router fiber.Router, h *AuthHandler) {
	auth := router.Group("/auth")
	auth.Post("/login", h.Login)
	auth.Post("/logout", RequireUserID, h.Logout)
	auth.Post("/otp/send", h.SendRegistrationOTP)
	auth.Post("/otp/verify", h.VerifyRegistrationOTP)
	auth.Post("/user/register", h.Register(domain.AccountTypeUser))
	auth.Post("/merchant/register", h.Register(domain.AccountTypeMerchant))

	internal := router.Group("/internal")
	internal.Post("/auth/verify-access-token", h.VerifyAccessToken)
	internal.Post("/auth/verify-transaction-token", h.VerifyTransactionToken)
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