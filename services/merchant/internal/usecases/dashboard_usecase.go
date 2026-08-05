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
	VPA          string
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
	qrRepo       domain.QRRepository
}

func NewDashboardUseCase(sr domain.MerchantRepository, rr domain.RewardRepository, cr domain.CampaignRepository, qr domain.QRRepository) *DashboardUseCase {
	return &DashboardUseCase{
		storeRepo:    sr,
		rewardRepo:   rr,
		campaignRepo: cr,
		qrRepo:       qr,
	}
}

// GetDashboard retrieves aggregated metrics, recent scan history, and active promotional campaigns for a merchant store.
func (uc *DashboardUseCase) GetDashboard(ctx context.Context, id string) (*DashboardOutput, error) {
	var storeID uuid.UUID
	var merchantIDStr string
	var vpaStr string
	var businessName string

	if id != "" {
		if parsedUUID, err := uuid.Parse(id); err == nil {
			if store, err := uc.storeRepo.GetStoreProfileByStoreID(ctx, parsedUUID); err == nil && store != nil && store.ID != uuid.Nil {
				storeID = store.ID
				merchantIDStr = store.DisplayID
				vpaStr = store.VPA
				businessName = store.LegalName
				if store.DBAName != nil && *store.DBAName != "" {
					businessName = *store.DBAName
				}
			} else if storeByOwner, err := uc.storeRepo.GetStoreProfileByOwnerID(ctx, parsedUUID); err == nil && storeByOwner != nil && storeByOwner.ID != uuid.Nil {
				storeID = storeByOwner.ID
				merchantIDStr = storeByOwner.DisplayID
				vpaStr = storeByOwner.VPA
				businessName = storeByOwner.LegalName
				if storeByOwner.DBAName != nil && *storeByOwner.DBAName != "" {
					businessName = *storeByOwner.DBAName
				}
			} else {
				storeID = parsedUUID
				merchantIDStr = "MM-8823-XA"
				vpaStr = "guest@moneymate"
				businessName = "Guest Merchant"
			}
		}
	}

	if storeID == uuid.Nil {
		storeID = uuid.New()
		merchantIDStr = "MM-9823-XA"
		vpaStr = "guest@moneymate"
		businessName = "Guest Merchant"
	}

	summary, err := uc.rewardRepo.GetSummaryByStoreID(ctx, storeID)
	if err != nil || summary == nil {
		summary = &domain.RewardSummary{
			AvailableBalance:       0.00,
			TotalScans:             0,
			PremiumPoints:          0,
			WeeklyGrowthPercentage: 0.00,
		}
	}

	txs, err := uc.qrRepo.GetQRTransactionsByStoreID(ctx, storeID, 10, 0)
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
			if len(tx.CustomerDisplayID) > 1 {
				initial = tx.CustomerDisplayID[:2]
			}
			dashTxs = append(dashTxs, DashboardTx{
				Time:     tx.CreatedAt.Format("Today, 15:04 PM"),
				Customer: tx.CustomerDisplayID,
				Initial:  initial,
				Color:    color,
				Amount:   fmt.Sprintf("$%.2f", tx.BillAmount),
				Reward:   fmt.Sprintf("$%.2f", tx.RewardIssued),
				Status:   "Settled",
			})
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
	}

	todayCount, _ := uc.qrRepo.GetTodayQRScanCount(ctx, storeID)
	todayVolume, _ := uc.qrRepo.GetTodayQRScanVolume(ctx, storeID)

	stats := []DashboardStat{
		{
			Title:          "Total Rewards Issued",
			Value:          fmt.Sprintf("$%.2f", summary.AvailableBalance),
			Icon:           "workspace_premium",
			IconColorClass: "text-primary bg-primary/10",
			BorderClass:    "border-l-primary",
			TrendText:      fmt.Sprintf("+%.0f%% this week", summary.WeeklyGrowthPercentage),
			TrendType:      "neutral",
		},
		{
			Title:          "Total QR Scan Volume",
			Value:          fmt.Sprintf("$%.2f", todayVolume),
			Icon:           "payments",
			IconColorClass: "text-tertiary bg-tertiary/10",
			BorderClass:    "border-l-tertiary",
			TrendText:      "Today's scan volume",
			TrendType:      "neutral",
		},
		{
			Title:          "Active Campaigns",
			Value:          fmt.Sprintf("%d Active", activeCount),
			Icon:           "local_offer",
			IconColorClass: "text-secondary bg-secondary-fixed/35",
			BorderClass:    "border-l-secondary",
			TrendText:      "Currently running",
			TrendType:      "neutral",
		},
		{
			Title:          "Customers Rewarded",
			Value:          fmt.Sprintf("%d", todayCount),
			Icon:           "person",
			IconColorClass: "text-outline bg-outline/10",
			BorderClass:    "border-l-outline",
			TrendText:      "Today's total scans",
			TrendType:      "neutral",
		},
	}

	return &DashboardOutput{
		Stats:        stats,
		Transactions: dashTxs,
		Campaigns:    dashCamps,
		Balance:      summary.AvailableBalance,
		MerchantID:   merchantIDStr,
		VPA:          vpaStr,
		BusinessName: businessName,
	}, nil
}
