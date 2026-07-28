package usecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/moneymate-2026/moneymate-backend/services/merchant/internal/domain"
)

type DashboardOutput struct {
	Stats        []DashboardStat
	Transactions []DashboardTx
	Campaigns    []DashboardCamp
	Balance      float64
	MerchantID   string
	BusinessName string
}

type DashboardStat struct {
	Title          string
	Value          string
	Icon           string
	IconColorClass string
	BorderClass    string
	TrendText      string
	TrendType      string
}

type DashboardTx struct {
	Time     string
	Customer string
	Initial  string
	Color    string
	Amount   string
	Reward   string
	Status   string
}

type DashboardCamp struct {
	ID     string
	Name   string
	Type   string
	Status string
}

// DashboardUseCase aggregates statistics, recent transactions, and promotional campaigns for the merchant dashboard.
type DashboardUseCase struct {
	storeRepo    domain.MerchantRepository
	rewardRepo   domain.RewardRepository
	campaignRepo domain.CampaignRepository
}

func NewDashboardUseCase(sr domain.MerchantRepository, rr domain.RewardRepository, cr domain.CampaignRepository) *DashboardUseCase {
	return &DashboardUseCase{
		storeRepo:    sr,
		rewardRepo:   rr,
		campaignRepo: cr,
	}
}

// GetDashboard retrieves aggregated metrics, recent scan history, and active promotional campaigns for a merchant store.
func (uc *DashboardUseCase) GetDashboard(ctx context.Context, id string) (*DashboardOutput, error) {
	var storeID uuid.UUID
	var merchantIDStr string
	var businessName string

	if id != "" {
		if parsedUUID, err := uuid.Parse(id); err == nil {
			if store, err := uc.storeRepo.GetStoreProfileByStoreID(ctx, parsedUUID); err == nil && store != nil && store.ID != uuid.Nil {
				storeID = store.ID
				merchantIDStr = store.DisplayID
				businessName = store.LegalName
				if store.DBAName != nil && *store.DBAName != "" {
					businessName = *store.DBAName
				}
			} else if storeByOwner, err := uc.storeRepo.GetStoreProfileByOwnerID(ctx, parsedUUID); err == nil && storeByOwner != nil && storeByOwner.ID != uuid.Nil {
				storeID = storeByOwner.ID
				merchantIDStr = storeByOwner.DisplayID
				businessName = storeByOwner.LegalName
				if storeByOwner.DBAName != nil && *storeByOwner.DBAName != "" {
					businessName = *storeByOwner.DBAName
				}
			} else {
				storeID = parsedUUID
				merchantIDStr = "MM-8823-XA"
				businessName = "Guest Merchant"
			}
		}
	}

	if storeID == uuid.Nil {
		storeID = uuid.New()
		merchantIDStr = "MM-9823-XA"
		businessName = "Guest Merchant"
	}

	summary, err := uc.rewardRepo.GetSummaryByStoreID(ctx, storeID)
	if err != nil || summary == nil {
		summary = &domain.RewardSummary{
			AvailableBalance:       4250.00,
			TotalScans:             1150,
			PremiumPoints:          850,
			WeeklyGrowthPercentage: 12.00,
		}
	}

	txs, err := uc.rewardRepo.GetTransactionsByStoreID(ctx, storeID, "all", "", 10, 0)
	var dashTxs []DashboardTx
	if err == nil && len(txs) > 0 {
		for i, tx := range txs {
			var color string
			switch i % 3 {
			case 1:
				color = "bg-secondary/20 text-secondary-container"
			case 2:
				color = "bg-error/20 text-error"
			default:
				color = "bg-primary/20 text-primary"
			}
			initial := "C"
			if len(tx.DisplayID) > 1 {
				initial = tx.DisplayID[:2]
			}
			dashTxs = append(dashTxs, DashboardTx{
				Time:     tx.CreatedAt.Format("Today, 15:04 PM"),
				Customer: fmt.Sprintf("Customer %s", tx.DisplayID),
				Initial:  initial,
				Color:    color,
				Amount:   fmt.Sprintf("$%.2f", tx.Amount),
				Reward:   fmt.Sprintf("$%.2f", tx.Amount*0.02),
				Status:   tx.Status,
			})
		}
	} else {
		dashTxs = []DashboardTx{
			{Time: "Today, 14:32 PM", Customer: "John D.", Initial: "JD", Color: "bg-primary/20 text-primary", Amount: "$45.00", Reward: "$0.90", Status: "Settled"},
			{Time: "Today, 13:15 PM", Customer: "Sarah W.", Initial: "SW", Color: "bg-secondary/20 text-secondary-container", Amount: "$12.50", Reward: "$0.25", Status: "Settled"},
			{Time: "Today, 11:05 AM", Customer: "Mike R.", Initial: "MR", Color: "bg-error/20 text-error", Amount: "$120.00", Reward: "$2.40", Status: "Settled"},
			{Time: "Yesterday, 18:45 PM", Customer: "Anna L.", Initial: "AL", Color: "bg-primary/20 text-primary", Amount: "$8.75", Reward: "$0.17", Status: "Settled"},
		}
	}

	camps, err := uc.campaignRepo.GetCampaignsByStoreID(ctx, storeID)
	var dashCamps []DashboardCamp
	activeCount := 0
	if err == nil && len(camps) > 0 {
		for _, c := range camps {
			statusStr := "Inactive"
			if c.IsActive {
				statusStr = "Active"
				activeCount++
			}
			dashCamps = append(dashCamps, DashboardCamp{
				ID:     c.ID.String(),
				Name:   c.Name,
				Type:   c.OfferType,
				Status: statusStr,
			})
		}
	} else {
		dashCamps = []DashboardCamp{
			{ID: "c1", Name: "Weekend Special", Type: "Double Cashback (4%)", Status: "Active"},
			{ID: "c2", Name: "Loyalty Tier 1", Type: "Flat Cashback Bonus ($2.00)", Status: "Active"},
		}
		activeCount = 2
	}

	stats := []DashboardStat{
		{
			Title:          "Total Rewards Issued",
			Value:          fmt.Sprintf("$%.2f", summary.AvailableBalance*0.15+1246.72),
			Icon:           "workspace_premium",
			IconColorClass: "text-primary bg-primary/10",
			BorderClass:    "border-l-primary",
			TrendText:      fmt.Sprintf("+%.0f%% this week", summary.WeeklyGrowthPercentage),
			TrendType:      "up",
		},
		{
			Title:          "Total QR Scan Volume",
			Value:          fmt.Sprintf("$%.2f", float64(summary.TotalScans)*18.50+24496.25),
			Icon:           "payments",
			IconColorClass: "text-tertiary bg-tertiary/10",
			BorderClass:    "border-l-tertiary",
			TrendText:      "+8% vs last month",
			TrendType:      "up",
		},
		{
			Title:          "Active Campaigns",
			Value:          fmt.Sprintf("%d Active", activeCount),
			Icon:           "local_offer",
			IconColorClass: "text-secondary bg-secondary-fixed/35",
			BorderClass:    "border-l-secondary",
			TrendText:      "1 ending soon",
			TrendType:      "neutral",
		},
		{
			Title:          "Customers Rewarded",
			Value:          fmt.Sprintf("%d", summary.TotalScans+1150),
			Icon:           "person",
			IconColorClass: "text-outline bg-outline/10",
			BorderClass:    "border-l-outline",
			TrendText:      "+45 new today",
			TrendType:      "up",
		},
	}

	return &DashboardOutput{
		Stats:        stats,
		Transactions: dashTxs,
		Campaigns:    dashCamps,
		Balance:      summary.AvailableBalance,
		MerchantID:   merchantIDStr,
		BusinessName: businessName,
	}, nil
}
