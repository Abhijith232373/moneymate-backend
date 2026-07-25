package usecase

import (
	"time"

	"github.com/google/uuid"
	"github.com/moneymate-2026/moneymate-backend/auth/internal/domain"
)

// ── Register ──────────────────────────────────────────────────────

type RegisterRequest struct {
	Email       string
	Phone       string
	FullName    string
	Password    string
	AccountType domain.AccountType 
}
type RegisterResponse struct {
	UserID uuid.UUID
	Email  string
	Handle string
	Status string
}

// ── Login ─────────────────────────────────────────────────────────

type LoginRequest struct {
	Identifier string
	Password   string
	DeviceID   string
	UserAgent  string
	IPAddress  string
}

type LoginResponse struct {
	AccessToken      string
	RefreshToken     string
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
	User             UserSummary
}

type UserSummary struct {
	ID              uuid.UUID
	Email           string
	Handle          string
	FullName        string
	Status          string
	IsEmailVerified bool
}

// ── Logout ────────────────────────────────────────────────────────

type LogoutRequest struct {
	UserID       uuid.UUID
	RefreshToken string
	AllDevices   bool
}

// ── Registration OTP ────────────────────────────────────────────

type SendRegistrationOTPRequest struct {
	Email string
}
type SendRegistrationOTPResponse struct {
    Email             string `json:"email"`
    ExpiresIn         int    `json:"expires_in"`         
    ResendCooldownIn  int    `json:"resend_cooldown_in"`  
    MaxVerifyAttempts int    `json:"max_verify_attempts"` 
}
type VerifyRegistrationOTPRequest struct {
	Email string
	Code  string
}

type VerifyRegistrationOTPResponse struct {
	Email    string
	Verified bool
}

type RefreshTokenRequest struct {
    RefreshToken string
}


//admin

type AdminUserSummary struct {
	ID              string `json:"id"`
	Email           string `json:"email"`
	Phone           string `json:"phone,omitempty"`
	FullName        string `json:"full_name"`
	Handle          string `json:"handle"`
	Status          string `json:"status"`
	IsEmailVerified bool   `json:"is_email_verified"`
	IsPhoneVerified bool   `json:"is_phone_verified"`
	CreatedAt       string `json:"created_at"`
}

type UserDetail struct {
	AdminUserSummary
	Roles []string `json:"roles"`
}


type ListUsersRequest struct {
	Status   string
	Search   string
	SortBy   string
	SortDesc bool
	Page     int
	PageSize int
}

type ListUsersResponse struct {
	Users      []AdminUserSummary `json:"users"`
	TotalCount int64         `json:"total_count"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
}

type UpdateUserRequest struct {
	FullName *string
	Email    *string
	Phone    *string
	Password *string 
}

// ── Role Management ──────────────────────────────────────────────

type CreateRoleRequest struct {
	Name        string
	Description *string
}

type UpdateRoleRequest struct {
	Name        *string
	Description *string
}

type RoleSummary struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type AssignRoleRequest struct {
	UserID    uuid.UUID
	RoleID    uuid.UUID
	GrantedBy *uuid.UUID 
}
type RefreshTokenResponse struct {
    AccessToken     string    `json:"access_token"`
    RefreshToken    string    `json:"refresh_token"`
    AccessExpiresAt time.Time `json:"access_expires_at"`
}