package app

import (
	"context"
	"fmt"
	"log"

	"github.com/moneymate-2026/moneymate-backend/auth/config"
	"github.com/moneymate-2026/moneymate-backend/auth/internal/domain"
	"github.com/moneymate-2026/moneymate-backend/auth/internal/infra/hasher"
	"github.com/moneymate-2026/moneymate-backend/auth/internal/infra/idgen"
)



func seedAdmin(ctx context.Context, userRepo domain.UserRepository, roleRepo domain.RoleRepository, h *hasher.Argon2Hasher, g *idgen.Generator, cfg *config.Config) error {    if cfg.Admin.Email == "" {
        log.Println("[SEED] no ADMIN_EMAIL set, skipping admin seed")
        return nil
    }

    exists, err := userRepo.EmailExists(ctx, cfg.Admin.Email)
    if err != nil {
        return fmt.Errorf("check admin exists: %w", err)
    }
    if exists {
        log.Println("[SEED] admin already exists, skipping")
        return nil
    }

    passwordHash, err := h.Hash(cfg.Admin.Password)
    if err != nil {
        return fmt.Errorf("hash admin password: %w", err)
    }

    userID, err := g.NewV7()
    if err != nil {
        return fmt.Errorf("generate admin id: %w", err)
    }

    role, err := roleRepo.GetByName(ctx, "admin")
    if err != nil {
        return fmt.Errorf("resolve admin role: %w", err)
    }

    user := &domain.User{
        ID:           userID,
        Email:        cfg.Admin.Email,
        FullName:     "System Admin",
        Handle:       "admin@moneymate",
        PasswordHash: &passwordHash,
        Status:       domain.UserStatusActive,
    }
    if err := userRepo.Create(ctx, user); err != nil {
        return fmt.Errorf("create admin user: %w", err)
    }
    if err := userRepo.VerifyEmail(ctx, user.ID); err != nil {
        return fmt.Errorf("verify admin email: %w", err)
    }
    if err := roleRepo.AssignRoleToUser(ctx, user.ID, role.ID, nil); err != nil {
        return fmt.Errorf("assign admin role: %w", err)
    }

    log.Printf("[SEED] admin created: %s", cfg.Admin.Email)
    return nil
}