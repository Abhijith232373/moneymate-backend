package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moneymate-2026/moneymate-backend/services/merchant/internal/domain"
	"github.com/moneymate-2026/moneymate-backend/services/merchant/sqlc/generated"
)

// SubscriptionRepo implements domain.SubscriptionRepository using PostgreSQL via pgxpool and SQLC generated queries.
// It is designed to handle dynamic pricing catalogs, atomic tier upgrades, and immutable billing ledgers for millions of users.
type SubscriptionRepo struct {
	// db holds the PostgreSQL connection pool.
	db *pgxpool.Pool
	// queries holds the sqlc-generated type-safe database methods.
	queries *generated.Queries
}

// NewSubscriptionRepo initializes and returns a new SubscriptionRepo instance.
func NewSubscriptionRepo(db *pgxpool.Pool) domain.SubscriptionRepository {
	return &SubscriptionRepo{
		db:      db,
		queries: generated.New(db),
	}
}

// GetAvailablePlans retrieves all active pricing tiers from the catalog, automatically seeding production tiers
// if the catalog is uninitialized, and dynamically flags which plan is currently assigned to the querying store.
func (r *SubscriptionRepo) GetAvailablePlans(ctx context.Context, storeID uuid.UUID) ([]*domain.SubscriptionPlanDetail, error) {
	// 1. Ensure catalog is seeded with standard production tiers matching UI specifications
	if err := r.seedCatalogIfEmpty(ctx); err != nil {
		return nil, fmt.Errorf("seed subscription catalog: %w", err)
	}

	// 2. Determine the store's current active plan code
	currentPlanCode := "essential"
	if storeID != uuid.Nil {
		sub, err := r.GetStoreSubscription(ctx, storeID)
		if err == nil && sub != nil {
			currentPlanCode = sub.PlanCode
		}
	}

	rows, err := r.queries.GetSubscriptionPlans(ctx)
	if err != nil {
		return nil, fmt.Errorf("query subscription plans: %w", err)
	}

	var plans []*domain.SubscriptionPlanDetail
	for _, row := range rows {
		var features []domain.PlanFeature
		if len(row.Features) > 0 {
			_ = json.Unmarshal(row.Features, &features)
		}
		if features == nil {
			features = []domain.PlanFeature{}
		}

		priceVal, _ := row.Price.Float64Value()

		p := &domain.SubscriptionPlanDetail{
			ID:                 row.ID,
			PlanCode:           row.PlanCode,
			Name:               row.Name,
			Price:              priceVal.Float64,
			BillingCycle:       string(row.BillingCycle),
			Description:        row.Description,
			MaxActiveCampaigns: int(row.MaxActiveCampaigns),
			IsMostPopular:      row.IsMostPopular,
			Features:           features,
			IsCurrent:          (row.PlanCode == currentPlanCode),
		}
		plans = append(plans, p)
	}

	return plans, nil
}

// seedCatalogIfEmpty checks if pricing tiers exist; if not, it populates the database with Essential, Growth, and Enterprise plans.
func (r *SubscriptionRepo) seedCatalogIfEmpty(ctx context.Context) error {
	existing, err := r.queries.GetSubscriptionPlans(ctx)
	if err == nil && len(existing) > 0 {
		return nil
	}

	essentialFeatures, _ := json.Marshal([]domain.PlanFeature{
		{Name: "Standard QR Payments", Included: true},
		{Name: "1 Active Offer Campaign", Included: true},
		{Name: "Standard Dashboard Analytics", Included: true},
		{Name: "Support via email", Included: true},
		{Name: "Custom API Integrations", Included: false},
		{Name: "Dedicated Account Manager", Included: false},
	})

	growthFeatures, _ := json.Marshal([]domain.PlanFeature{
		{Name: "Unlimited QR Payments", Included: true},
		{Name: "5 Active Offer Campaigns", Included: true},
		{Name: "Advanced Analytics & Insights", Included: true},
		{Name: "SMS Customer Notifications", Included: true},
		{Name: "Priority Email & Chat Support", Included: true},
		{Name: "Custom API Integrations", Included: false},
	})

	enterpriseFeatures, _ := json.Marshal([]domain.PlanFeature{
		{Name: "Everything in Growth", Included: true},
		{Name: "Unlimited Active Campaigns", Included: true},
		{Name: "Custom API & Webhook Integrations", Included: true},
		{Name: "Dedicated Account Manager", Included: true},
		{Name: "24/7 Phone & Slack Support", Included: true},
	})

	var p0, p29, p99 pgtype.Numeric
	_ = p0.Scan("0.00")
	_ = p29.Scan("29.00")
	_ = p99.Scan("99.00")

	_, _ = r.queries.CreateSubscriptionPlan(ctx, generated.CreateSubscriptionPlanParams{
		PlanCode:           "essential",
		Name:               "Essential",
		Price:              p0,
		BillingCycle:       generated.BillingCycleTypeMonthly,
		Description:        "For new merchants starting their loyalty and scan payments journey.",
		MaxActiveCampaigns: 1,
		IsMostPopular:      false,
		Features:           essentialFeatures,
		IsActive:           true,
	})

	_, _ = r.queries.CreateSubscriptionPlan(ctx, generated.CreateSubscriptionPlanParams{
		PlanCode:           "growth",
		Name:               "Growth",
		Price:              p29,
		BillingCycle:       generated.BillingCycleTypeMonthly,
		Description:        "For growing businesses wanting to supercharge sales and customer retention.",
		MaxActiveCampaigns: 5,
		IsMostPopular:      true,
		Features:           growthFeatures,
		IsActive:           true,
	})

	_, _ = r.queries.CreateSubscriptionPlan(ctx, generated.CreateSubscriptionPlanParams{
		PlanCode:           "enterprise",
		Name:               "Enterprise",
		Price:              p99,
		BillingCycle:       generated.BillingCycleTypeMonthly,
		Description:        "For large retail chains needing bespoke integrations and dedicated support.",
		MaxActiveCampaigns: -1,
		IsMostPopular:      false,
		Features:           enterpriseFeatures,
		IsActive:           true,
	})

	return nil
}

// GetStoreSubscription retrieves a merchant's active billing record. If uninitialized, it inspects the stores table
// and initializes a corresponding merchant subscription record to guarantee seamless onboarding.
func (r *SubscriptionRepo) GetStoreSubscription(ctx context.Context, storeID uuid.UUID) (*domain.MerchantSubscription, error) {
	row, err := r.queries.GetSubscriptionByStoreID(ctx, storeID)
	if errors.Is(err, pgx.ErrNoRows) {
		return r.initializeStoreSubscription(ctx, storeID)
	}
	if err != nil {
		return nil, fmt.Errorf("query store subscription: %w", err)
	}

	return &domain.MerchantSubscription{
		ID:                 row.ID,
		StoreID:            row.StoreID,
		PlanCode:           row.PlanCode,
		Status:             string(row.Status),
		BillingCycle:       string(row.BillingCycle),
		CurrentPeriodStart: row.CurrentPeriodStart,
		CurrentPeriodEnd:   row.CurrentPeriodEnd,
		AutoRenew:          row.AutoRenew,
		CreatedAt:          row.CreatedAt,
		UpdatedAt:          row.UpdatedAt,
	}, nil
}

// initializeStoreSubscription creates a default active subscription for a store based on its current plan setting in the stores table.
func (r *SubscriptionRepo) initializeStoreSubscription(ctx context.Context, storeID uuid.UUID) (*domain.MerchantSubscription, error) {
	planCode := "essential"
	var storePlan string
	err := r.db.QueryRow(ctx, `SELECT plan::text FROM stores WHERE id = $1;`, storeID).Scan(&storePlan)
	if err == nil && storePlan != "" {
		planCode = storePlan
	}

	now := time.Now().UTC()
	row, err := r.queries.UpsertMerchantSubscription(ctx, generated.UpsertMerchantSubscriptionParams{
		ID:                 uuid.New(),
		StoreID:            storeID,
		PlanCode:           planCode,
		Status:             generated.SubscriptionBillingStatusActive,
		BillingCycle:       generated.BillingCycleTypeMonthly,
		CurrentPeriodStart: now,
		CurrentPeriodEnd:   now.AddDate(0, 0, 30),
		AutoRenew:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("insert initial merchant subscription: %w", err)
	}

	return &domain.MerchantSubscription{
		ID:                 row.ID,
		StoreID:            row.StoreID,
		PlanCode:           row.PlanCode,
		Status:             string(row.Status),
		BillingCycle:       string(row.BillingCycle),
		CurrentPeriodStart: row.CurrentPeriodStart,
		CurrentPeriodEnd:   row.CurrentPeriodEnd,
		AutoRenew:          row.AutoRenew,
		CreatedAt:          row.CreatedAt,
		UpdatedAt:          row.UpdatedAt,
	}, nil
}

// UpdateStorePlan executes an atomic database transaction that updates the merchant's subscription tier
// and synchronizes the core store record.
func (r *SubscriptionRepo) UpdateStorePlan(ctx context.Context, storeID uuid.UUID, newPlanCode string) (*domain.MerchantSubscription, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin subscription update tx: %w", err)
	}
	defer tx.Rollback(ctx)

	qTx := r.queries.WithTx(tx)

	// 1. Lock and retrieve current subscription
	var sub domain.MerchantSubscription
	lockQuery := `
		SELECT id, store_id, plan_code, status::text, billing_cycle::text, current_period_start, current_period_end, auto_renew, created_at, updated_at
		FROM merchant_subscriptions
		WHERE store_id = $1
		FOR UPDATE;
	`
	err = tx.QueryRow(ctx, lockQuery, storeID).Scan(
		&sub.ID, &sub.StoreID, &sub.PlanCode, &sub.Status, &sub.BillingCycle,
		&sub.CurrentPeriodStart, &sub.CurrentPeriodEnd, &sub.AutoRenew, &sub.CreatedAt, &sub.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		_, _ = r.initializeStoreSubscription(ctx, storeID)
		err = tx.QueryRow(ctx, lockQuery, storeID).Scan(
			&sub.ID, &sub.StoreID, &sub.PlanCode, &sub.Status, &sub.BillingCycle,
			&sub.CurrentPeriodStart, &sub.CurrentPeriodEnd, &sub.AutoRenew, &sub.CreatedAt, &sub.UpdatedAt,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("lock store subscription: %w", err)
	}

	oldPlanCode := sub.PlanCode
	if oldPlanCode == newPlanCode {
		return &sub, nil // Idempotent no-op if already on the requested tier
	}

	// 2. Update merchant_subscriptions table using sqlc generated query
	row, err := qTx.UpsertMerchantSubscription(ctx, generated.UpsertMerchantSubscriptionParams{
		ID:                 sub.ID,
		StoreID:            storeID,
		PlanCode:           newPlanCode,
		Status:             generated.SubscriptionBillingStatus(sub.Status),
		BillingCycle:       generated.BillingCycleType(sub.BillingCycle),
		CurrentPeriodStart: sub.CurrentPeriodStart,
		CurrentPeriodEnd:   sub.CurrentPeriodEnd,
		AutoRenew:          sub.AutoRenew,
	})
	if err != nil {
		return nil, fmt.Errorf("update merchant subscription: %w", err)
	}

	// 3. Synchronize core stores table plan enum column
	err = qTx.UpdateStorePlanEnum(ctx, generated.UpdateStorePlanEnumParams{
		ID:   storeID,
		Plan: generated.SubscriptionPlan(newPlanCode),
	})
	if err != nil {
		return nil, fmt.Errorf("sync store plan enum: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit subscription update tx: %w", err)
	}

	return &domain.MerchantSubscription{
		ID:                 row.ID,
		StoreID:            row.StoreID,
		PlanCode:           row.PlanCode,
		Status:             string(row.Status),
		BillingCycle:       string(row.BillingCycle),
		CurrentPeriodStart: row.CurrentPeriodStart,
		CurrentPeriodEnd:   row.CurrentPeriodEnd,
		AutoRenew:          row.AutoRenew,
		CreatedAt:          row.CreatedAt,
		UpdatedAt:          row.UpdatedAt,
	}, nil
}
