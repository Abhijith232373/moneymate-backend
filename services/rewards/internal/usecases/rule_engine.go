package usecases

import "github.com/moneymate-2026/moneymate-backend/services/rewards/internal/domain"

type RollBPSFunc func(minBPS, maxBPS int32) int32

type RewardCalculation struct {
	Eligible          bool
	PercentageBPS     int32
	RewardAmountPaise int64
	Reason            string
}

func CalculateReward(amountPaise int64, rule domain.RewardRule, roll RollBPSFunc) RewardCalculation {
	if !rule.Active {
		return RewardCalculation{Eligible: false, Reason: "rule inactive"}
	}
	if amountPaise <= 0 {
		return RewardCalculation{Eligible: false, Reason: "invalid amount"}
	}
	if amountPaise < rule.MinTransactionAmountPaise {
		return RewardCalculation{Eligible: false, Reason: "below threshold"}
	}

	bps := roll(rule.MinPercentageBPS, rule.MaxPercentageBPS)
	if bps < rule.MinPercentageBPS {
		bps = rule.MinPercentageBPS
	}
	if bps > rule.MaxPercentageBPS {
		bps = rule.MaxPercentageBPS
	}

	reward := amountPaise * int64(bps) / 10000
	if reward > rule.MaxPayoutAmountPaise {
		reward = rule.MaxPayoutAmountPaise
	}
	if reward <= 0 {
		return RewardCalculation{Eligible: false, PercentageBPS: bps, Reason: "zero reward"}
	}

	return RewardCalculation{
		Eligible:          true,
		PercentageBPS:     bps,
		RewardAmountPaise: reward,
	}
}
