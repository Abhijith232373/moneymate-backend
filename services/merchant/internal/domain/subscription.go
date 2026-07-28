package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// PlanFeature represents an individual feature checklist item displayed on the pricing card UI.
type PlanFeature struct {
	Name     string
	Included bool
}

// SubscriptionPlanDetail represents a complete pricing tier catalog item with dynamic feature availability.
type SubscriptionPlanDetail struct {
	ID                 uuid.UUID
	PlanCode           string
	Name               string
	Price              float64
	BillingCycle       string
	Description        string
	MaxActiveCampaigns int
	IsMostPopular      bool
	Features           []PlanFeature
	IsCurrent          bool
}

// MerchantSubscription represents the active billing state and plan assignment for a merchant store.
type MerchantSubscription struct {
	ID                 uuid.UUID
	StoreID            uuid.UUID
	PlanCode           string
	Status             string
	BillingCycle       string
	CurrentPeriodStart time.Time
	CurrentPeriodEnd   time.Time
	AutoRenew          bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// SubscriptionChangeLog represents an immutable audit ledger entry for tier upgrades and downgrades.
type SubscriptionChangeLog struct {
	ID           uuid.UUID
	StoreID      uuid.UUID
	OldPlanCode  string
	NewPlanCode  string
	ChangeReason string
	ChangedAt    time.Time
}

// SubscriptionRepository defines the strict data access contract for subscription plans, billing states, and audit logs.
type SubscriptionRepository interface {
	// GetAvailablePlans retrieves all active subscription tiers from the catalog, marking the store's current plan dynamically.
	GetAvailablePlans(ctx context.Context, storeID uuid.UUID) ([]*SubscriptionPlanDetail, error)
	// GetStoreSubscription fetches the active subscription record for a specific merchant store, bootstrapping default Essential if missing.
	GetStoreSubscription(ctx context.Context, storeID uuid.UUID) (*MerchantSubscription, error)
	// UpdateStorePlan atomically upgrades or downgrades a store's subscription tier, updates core store records, and logs an immutable audit trail.
	UpdateStorePlan(ctx context.Context, storeID uuid.UUID, newPlanCode string, reason string) (*MerchantSubscription, error)
}
