package http

type EarningsResponse struct {
	TotalScans          int64           `json:"total_scans"`
	TotalEarned         float64         `json:"total_earned"`
	FormattedTotal      string          `json:"formatted_total"`
	RequestedMilestones map[int32]bool `json:"requested_milestones"`
}

type RequestPayoutRequest struct {
	MilestoneScans int32   `json:"milestone_scans"`
	RewardAmount   float64 `json:"reward_amount"`
}
