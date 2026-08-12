package http

type RegisterStoreRequest struct {

	OwnerName         string `json:"owner_name"`
	ContactEmail      string `json:"contact_email"`
	MobileNumber      string `json:"mobile_number"`
	LegalName         string `json:"legal_name"`
	DBAName           string `json:"dba_name,omitempty"`
	BusinessType      string `json:"business_type"`
	TaxID             string `json:"tax_id,omitempty"`
	RegisteredAddress string `json:"registered_address"`
	AadhaarNumber     string `json:"aadhaar_number"`
	AadhaarDocURL     string `json:"aadhaar_doc_url"`
	ShopLicenseURL    string `json:"shop_license_url"`
	Password          string `json:"password"`
	ConfirmPassword   string `json:"confirm_password"`
}

type RegisterStoreResponse struct {
	StoreID      string `json:"store_id"`
	DisplayID    string `json:"display_id"`
	VPA          string `json:"vpa"`
	QRCodeBase64 string `json:"qr_code_base64"`
	Status       string `json:"status"`
	Plan         string `json:"plan"`
	Token        string `json:"token"`
}

type LoginStoreRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginStoreResponse struct {
	StoreID      string `json:"store_id"`

	DisplayID    string `json:"display_id"`
	VPA          string `json:"vpa"`
	LegalName    string `json:"legal_name"`
	Status       string `json:"status"`
	Plan         string `json:"plan"`
	Token        string `json:"token"`
}

type GetStoreResponse struct {
	StoreID   string `json:"store_id"`
	DisplayID string `json:"display_id"`
	VPA          string `json:"vpa"`
	QRCodeBase64 string `json:"qr_code_base64"`
	Status       string `json:"status"`
	Plan      string `json:"plan"`
	LegalName string `json:"legal_name"`
}

type CreateCampaignRequest struct {
	Name            string  `json:"name"`
	RedeemCode      string  `json:"redeem_code"`
	OfferCategory   string  `json:"offer_category"`
	OfferType       string  `json:"offer_type"`
	RewardValue     float64 `json:"reward_value"`
	MinBillAmount   float64 `json:"min_bill_amount"`
	RedemptionLimit int32   `json:"redemption_limit"`
	TargetAudience  string  `json:"target_audience"`
	StartDate       string  `json:"start_date"` // YYYY-MM-DDTHH:mm:ss
	EndDate         string  `json:"end_date"`   // YYYY-MM-DDTHH:mm:ss
	BannerURL       string  `json:"banner_url,omitempty"`
}

type CampaignResponse struct {
	ID              string  `json:"id"`
	StoreID         string  `json:"store_id"`
	Name            string  `json:"name"`
	RedeemCode      string  `json:"redeem_code"`
	OfferCategory   string  `json:"offer_category"`
	OfferType       string  `json:"offer_type"`
	RewardValue     float64 `json:"reward_value"`
	MinBillAmount   float64 `json:"min_bill_amount"`
	RedemptionLimit int32   `json:"redemption_limit"`
	StartDate       string  `json:"start_date"`
	EndDate         string  `json:"end_date"`
	BannerURL       string  `json:"banner_url,omitempty"`
	Status          string  `json:"status"`
	CreatedAt       string  `json:"created_at"`
}

// RewardSummaryResponse represents the JSON payload sent to the frontend Rewards Center dashboard.
type RewardSummaryResponse struct {
	StoreID                string  `json:"store_id"`
	AvailableBalance       float64 `json:"available_balance"`
	TotalScans             int64   `json:"total_scans"`
	PremiumPoints          int64   `json:"premium_points"`
	WeeklyGrowthPercentage float64 `json:"weekly_growth_percentage"`
	FormattedBalance       string  `json:"formatted_balance"`
	FormattedGrowth        string  `json:"formatted_growth"`
}

// RewardTransactionResponse represents a single item in the Rewards Center transaction history table.
type RewardTransactionResponse struct {
	ID              string  `json:"id"`
	StoreID         string  `json:"store_id"`
	CampaignName    string  `json:"campaign_name"`
	DisplayID       string  `json:"display_id"`
	Status          string  `json:"status"`
	Amount          float64 `json:"amount"`
	FormattedAmount string  `json:"formatted_amount"`
	TransactionType string  `json:"transaction_type"`
	CreatedAt       string  `json:"created_at"`
	FormattedDate   string  `json:"formatted_date"`
}

// RedeemBalanceRequest represents the incoming JSON payload when a merchant initiates a bank withdrawal.
type RedeemBalanceRequest struct {
	Amount                           float64 `json:"amount"`
	ConfirmBankTransferAuthorization bool    `json:"confirm_bank_transfer_authorization"`
}

// RedeemBalanceResponse represents the acknowledgment payload returned after successfully initiating a bank withdrawal.
type RedeemBalanceResponse struct {
	RedemptionID              string  `json:"redemption_id"`
	StoreID                   string  `json:"store_id"`
	AmountRedeemed            float64 `json:"amount_redeemed"`
	FormattedAmountRedeemed   string  `json:"formatted_amount_redeemed"`
	RemainingBalance          float64 `json:"remaining_balance"`
	FormattedRemainingBalance string  `json:"formatted_remaining_balance"`
	Status                    string  `json:"status"`
	ReferenceID               string  `json:"reference_id,omitempty"`
	Message                   string  `json:"message"`
}

// PlanFeatureResponse represents a single check or cross item in the UI feature checklist.
type PlanFeatureResponse struct {
	Name     string `json:"name"`
	Included bool   `json:"included"`
}

// SubscriptionPlanResponse represents the JSON payload for a pricing card on the "Choose Your Plan" dashboard.
type SubscriptionPlanResponse struct {
	ID                 string                `json:"id"`
	PlanCode           string                `json:"plan_code"`
	Name               string                `json:"name"`
	Price              float64               `json:"price"`
	FormattedPrice     string                `json:"formatted_price"`
	BillingCycle       string                `json:"billing_cycle"`
	FormattedCycle     string                `json:"formatted_cycle"`
	Description        string                `json:"description"`
	MaxActiveCampaigns int                   `json:"max_active_campaigns"`
	IsMostPopular      bool                  `json:"is_most_popular"`
	Features           []PlanFeatureResponse `json:"features"`
	IsCurrent          bool                  `json:"is_current"`
}

// CurrentSubscriptionResponse represents the active billing standing and renewal timeline for a merchant store.
type CurrentSubscriptionResponse struct {
	ID                 string `json:"id"`
	StoreID            string `json:"store_id"`
	PlanCode           string `json:"plan_code"`
	PlanName           string `json:"plan_name"`
	Status             string `json:"status"`
	BillingCycle       string `json:"billing_cycle"`
	CurrentPeriodStart string `json:"current_period_start"`
	CurrentPeriodEnd   string `json:"current_period_end"`
	FormattedRenewalDate string `json:"formatted_renewal_date"`
	AutoRenew          bool   `json:"auto_renew"`
}

// ChangePlanRequest represents the incoming JSON body when upgrading or downgrading subscription tiers.
type ChangePlanRequest struct {
	PlanCode string `json:"plan_code"`
}

// ChangePlanResponse represents the acknowledgment payload after successfully transitioning subscription tiers.
type ChangePlanResponse struct {
	ID       string `json:"id"`
	StoreID  string `json:"store_id"`
	PlanCode string `json:"plan_code"`
	PlanName string `json:"plan_name"`
	Status   string `json:"status"`
	Message  string `json:"message"`
}

// KYCDocumentItemResponse represents an individual compliance document card displayed on the KYC Status dashboard.
type KYCDocumentItemResponse struct {
	Title           string `json:"title"`
	Status          string `json:"status"`
	SubmittedAt     string `json:"submitted_at"`
	FormattedStatus string `json:"formatted_status"`
	URL             string `json:"url"`
	DocType         string `json:"doc_type"`
}

// KYCStatusResponse represents the JSON payload for the "KYC Status" verification dashboard.
type KYCStatusResponse struct {
	StoreID       string                    `json:"store_id"`
	Status        string                    `json:"status"`
	IsVerified    bool                      `json:"is_verified"`
	Message       string                    `json:"message"`
	NextReviewDue string                    `json:"next_review_due"`
	Documents     []KYCDocumentItemResponse `json:"documents"`
}

// UpdateKYCDocumentsRequest represents the incoming JSON body when submitting new or updated compliance files.
type UpdateKYCDocumentsRequest struct {
	AadhaarNumber  string `json:"aadhaar_number,omitempty"`
	AadhaarDocURL  string `json:"aadhaar_doc_url,omitempty"`
	ShopLicenseURL string `json:"shop_license_url,omitempty"`
}

// UpdateKYCDocumentsResponse represents the acknowledgment payload after submitting updated compliance documentation.
type UpdateKYCDocumentsResponse struct {
	StoreID    string `json:"store_id"`
	Status     string `json:"status"`
	IsVerified bool   `json:"is_verified"`
	// Message confirms that the documents were received and are being processed.
	Message string `json:"message"`
}

// ProfileResponse represents the merchant business and primary contact profile data for the Settings/Profile page.
type ProfileResponse struct {
	StoreID      string `json:"storeId,omitempty"`

	BusinessName string `json:"businessName"`
	DBAName      string `json:"dbaName"`
	Address      string `json:"address"`
	BusinessType string `json:"businessType"`
	TaxID        string `json:"taxId"`
	OwnerName    string `json:"ownerName"`
	Email        string `json:"email"`
	Mobile       string `json:"mobile"`
	ProfileImage string `json:"profileImage"`
	Status       string `json:"status"`
	DisplayID    string `json:"displayId"`
	VPA          string `json:"vpa"`
	QRCodeBase64 string `json:"qr_code_base64"`
	Plan         string `json:"plan"`
	CreatedAt    string `json:"createdAt"`
}

// UpdateProfileRequest represents the incoming payload from the merchant Settings/Profile page form.
type UpdateProfileRequest struct {
	BusinessName string `json:"businessName"`
	DBAName      string `json:"dbaName"`
	Address      string `json:"address"`
	BusinessType string `json:"businessType"`
	TaxID        string `json:"taxId"`
	OwnerName    string `json:"ownerName"`
	Email        string `json:"email"`
	Mobile       string `json:"mobile"`
	ProfileImage string `json:"profileImage"`
}

// TrendResponse represents the trend indicator on a dashboard stat card.
type TrendResponse struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

// StatCardResponse represents a single KPI card on the merchant dashboard.
type StatCardResponse struct {
	Title          string        `json:"title"`
	Value          string        `json:"value"`
	Icon           string        `json:"icon"`
	IconColorClass string        `json:"iconColorClass"`
	BorderClass    string        `json:"borderClass"`
	Trend          TrendResponse `json:"trend"`
}

// DashboardTransactionResponse represents a recent transaction in the dashboard table.
type DashboardTransactionResponse struct {
	Time     string `json:"time"`
	Customer string `json:"customer"`
	Initial  string `json:"initial"`
	Color    string `json:"color"`
	Amount   string `json:"amount"`
	Reward   string `json:"reward"`
	Status   string `json:"status"`
}

// DashboardCampaignResponse represents an active campaign summary in the dashboard.
type DashboardCampaignResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

// DashboardResponse represents the complete payload for the merchant dashboard overview.
type DashboardResponse struct {
	Stats        []StatCardResponse             `json:"stats"`
	Transactions []DashboardTransactionResponse `json:"transactions"`
	Campaigns    []DashboardCampaignResponse    `json:"campaigns"`
	Balance      float64                        `json:"balance"`
	MerchantID   string                         `json:"merchant_id"`
	VPA          string                         `json:"vpa"`
	BusinessName string                         `json:"business_name"`
}
