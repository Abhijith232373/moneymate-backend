package http

type WalletOverviewResponse struct {
	AvailableBalance        float64 `json:"available_balance"`
	FormattedBalance        string  `json:"formatted_balance"`
	TotalEarnings           float64 `json:"total_earnings"`
	FormattedTotalEarnings  string  `json:"formatted_total_earnings"`
	TotalRedeemed           float64 `json:"total_redeemed"`
	FormattedTotalRedeemed  string  `json:"formatted_total_redeemed"`
}

type WalletTransactionResponse struct {
	ID              string  `json:"id"`
	TransactionID   string  `json:"transaction_id"`
	Title           string  `json:"title"`
	Subtitle        string  `json:"subtitle"`
	Amount          float64 `json:"amount"`
	FormattedAmount string  `json:"formatted_amount"`
	Date            string  `json:"date"`
	Time            string  `json:"time"`
}

type WalletResponse struct {
	Overview     WalletOverviewResponse      `json:"overview"`
	Transactions []WalletTransactionResponse `json:"transactions"`
}
