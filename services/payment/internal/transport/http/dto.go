package http

type createWalletRequest struct {
	// user_id removed — the wallet is always created for the authenticated
	// caller (see middleware.go). Kept as an empty struct placeholder in
	// case an admin-only "create wallet for user X" endpoint is added later
	// behind a separate, explicitly-privileged route.
}

type transferRequest struct {
	FromAccountID  string `json:"from_account_id" validate:"required"`
	ToAccountID    string `json:"to_account_id" validate:"required"`
	Amount         string `json:"amount" validate:"required"`
	IdempotencyKey string `json:"idempotency_key" validate:"required"`
	Description    string `json:"description"`
}

type walletResponse struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	Type     string `json:"type"`
	Currency string `json:"currency"`
	Balance  string `json:"balance"` 
}

type transactionResponse struct {
	ID             string `json:"id"`
	FromAccountID  string `json:"from_account_id"`
	ToAccountID    string `json:"to_account_id"`
	Amount         string `json:"amount"`
	Status         string `json:"status"`
	IdempotencyKey string `json:"idempotency_key"`
	Description    string `json:"description"`
	CreatedAt      string `json:"created_at"`
}

type transferResponse struct {
	Transaction transactionResponse `json:"transaction"`
	FromBalance string             `json:"from_balance"`
	ToBalance   string             `json:"to_balance"`
}

type createWalletInternalRequest struct {
	UserID string `json:"user_id" validate:"required"`
	Handle string `json:"handle" validate:"required"`
}


type initiateDepositRequest struct {
	AmountRupees float64 `json:"amount" validate:"required,gt=0"`
}

type confirmDepositRequest struct {
	OrderID   string `json:"razorpay_order_id" validate:"required"`
	PaymentID string `json:"razorpay_payment_id" validate:"required"`
	Signature string `json:"razorpay_signature" validate:"required"`
}

type requestWithdrawalRequest struct {
	AmountRupees   float64 `json:"amount" validate:"required,gt=0"`
	IdempotencyKey string  `json:"idempotency_key" validate:"required"`
}