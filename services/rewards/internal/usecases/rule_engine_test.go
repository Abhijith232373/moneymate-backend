package usecases

import (
	"testing"

	"github.com/moneymate-2026/moneymate-backend/services/rewards/internal/domain"
)

func TestCalculateReward(t *testing.T) {
	base := domain.RewardRule{
		MinPercentageBPS:          0,
		MaxPercentageBPS:          150,
		MinTransactionAmountPaise: 10000,
		MaxPayoutAmountPaise:      5000,
		Active:                    true,
	}

	tests := []struct {
		name       string
		amount     int64
		rule       domain.RewardRule
		rolledBPS  int32
		eligible   bool
		wantAmount int64
	}{
		{name: "below threshold", amount: 9999, rule: base, rolledBPS: 100, eligible: false},
		{name: "normal range", amount: 100000, rule: base, rolledBPS: 100, eligible: true, wantAmount: 1000},
		{name: "at cap", amount: 1000000, rule: base, rolledBPS: 150, eligible: true, wantAmount: 5000},
		{name: "inactive", amount: 100000, rule: func() domain.RewardRule { r := base; r.Active = false; return r }(), rolledBPS: 100, eligible: false},
		{name: "zero amount", amount: 0, rule: base, rolledBPS: 100, eligible: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateReward(tt.amount, tt.rule, func(_, _ int32) int32 { return tt.rolledBPS })
			if got.Eligible != tt.eligible {
				t.Fatalf("eligible = %v, want %v", got.Eligible, tt.eligible)
			}
			if got.RewardAmountPaise != tt.wantAmount {
				t.Fatalf("amount = %d, want %d", got.RewardAmountPaise, tt.wantAmount)
			}
		})
	}
}
