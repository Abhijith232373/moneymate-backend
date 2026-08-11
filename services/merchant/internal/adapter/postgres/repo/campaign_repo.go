package repo

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moneymate-2026/moneymate-backend/services/merchant/internal/domain"
	"github.com/moneymate-2026/moneymate-backend/services/merchant/sqlc/generated"
)

type CampaignRepo struct {
	db      *pgxpool.Pool
	queries *generated.Queries
}

func NewCampaignRepo(db *pgxpool.Pool) domain.CampaignRepository {
	return &CampaignRepo{
		db:      db,
		queries: generated.New(db),
	}
}

func (r *CampaignRepo) CreateCampaign(ctx context.Context, c *domain.Campaign) (*domain.Campaign, error) {
	var rewardValue pgtype.Numeric
	rewardValue.Scan(strconv.FormatFloat(c.RewardValue, 'f', -1, 64))

	var minBillAmount pgtype.Numeric
	minBillAmount.Scan(strconv.FormatFloat(c.MinBillAmount, 'f', -1, 64))

	params := generated.CreateCampaignParams{
		StoreID:         c.StoreID,
		Name:            c.Name,
		RedeemCode:      c.RedeemCode,
		OfferCategory:   c.OfferCategory,
		OfferType:       c.OfferType,
		RewardValue:     rewardValue,
		MinBillAmount:   minBillAmount,
		RedemptionLimit: c.RedemptionLimit,
		TargetAudience:  c.TargetAudience,
		BannerUrl:       c.BannerURL,
		StartDate:       c.StartDate,
		EndDate:         c.EndDate,
		Status:          c.Status,
	}

	res, err := r.queries.CreateCampaign(ctx, params)
	if err != nil {
		return nil, err
	}

	rv, _ := res.RewardValue.Float64Value()
	mb, _ := res.MinBillAmount.Float64Value()

	c.ID = res.ID
	c.StoreID = res.StoreID
	c.Name = res.Name
	c.RedeemCode = res.RedeemCode
	c.OfferCategory = res.OfferCategory
	c.OfferType = res.OfferType
	c.RewardValue = rv.Float64
	c.MinBillAmount = mb.Float64
	c.RedemptionLimit = res.RedemptionLimit
	c.TargetAudience = res.TargetAudience
	c.BannerURL = res.BannerUrl
	c.StartDate = res.StartDate
	c.EndDate = res.EndDate
	c.Status = res.Status
	c.CreatedAt = res.CreatedAt
	c.UpdatedAt = res.UpdatedAt

	return c, nil
}

func (r *CampaignRepo) GetCampaignsByStoreID(ctx context.Context, storeID uuid.UUID) ([]*domain.Campaign, error) {
	res, err := r.queries.GetCampaignsByStoreID(ctx, storeID)
	if err != nil {
		return nil, err
	}

	var campaigns []*domain.Campaign
	for _, c := range res {
		rv, _ := c.RewardValue.Float64Value()
		mb, _ := c.MinBillAmount.Float64Value()

		campaigns = append(campaigns, &domain.Campaign{
			ID:              c.ID,
			StoreID:         c.StoreID,
			Name:            c.Name,
			RedeemCode:      c.RedeemCode,
			OfferCategory:   c.OfferCategory,
			OfferType:       c.OfferType,
			RewardValue:     rv.Float64,
			MinBillAmount:   mb.Float64,
			RedemptionLimit: c.RedemptionLimit,
			TargetAudience:  c.TargetAudience,
			BannerURL:       c.BannerUrl,
			StartDate:       c.StartDate,
			EndDate:         c.EndDate,
			Status:          c.Status,
			CreatedAt:       c.CreatedAt,
			UpdatedAt:       c.UpdatedAt,
		})
	}
	return campaigns, nil
}

func (r *CampaignRepo) GetPublicCampaigns(ctx context.Context, limit, offset int) ([]*domain.Campaign, error) {
	if limit <= 0 {
		limit = 50
	}
	query := `
		SELECT 
			id, store_id, name, redeem_code, offer_category, offer_type, COALESCE(reward_value, 0), COALESCE(min_bill_amount, 0), redemption_limit, target_audience, banner_url, start_date, end_date, status, created_at, updated_at
		FROM campaigns
		WHERE status = 'active' AND end_date > NOW()
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2;`
	
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var campaigns []*domain.Campaign
	for rows.Next() {
		var c domain.Campaign
		if err := rows.Scan(
			&c.ID, &c.StoreID, &c.Name, &c.RedeemCode, &c.OfferCategory, &c.OfferType, &c.RewardValue, &c.MinBillAmount,
			&c.RedemptionLimit, &c.TargetAudience, &c.BannerURL, &c.StartDate, &c.EndDate, &c.Status, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		campaigns = append(campaigns, &c)
	}
	return campaigns, nil
}


func (r *CampaignRepo) GetCampaignByID(ctx context.Context, id uuid.UUID) (*domain.Campaign, error) {
	res, err := r.queries.GetCampaignByID(ctx, id)
	if err != nil {
		return nil, err
	}

	rv, _ := res.RewardValue.Float64Value()
	mb, _ := res.MinBillAmount.Float64Value()

	return &domain.Campaign{
		ID:              res.ID,
		StoreID:         res.StoreID,
		Name:            res.Name,
		RedeemCode:      res.RedeemCode,
		OfferCategory:   res.OfferCategory,
		OfferType:       res.OfferType,
		RewardValue:     rv.Float64,
		MinBillAmount:   mb.Float64,
		RedemptionLimit: res.RedemptionLimit,
		TargetAudience:  res.TargetAudience,
		BannerURL:       res.BannerUrl,
		StartDate:       res.StartDate,
		EndDate:         res.EndDate,
		Status:          res.Status,
		CreatedAt:       res.CreatedAt,
		UpdatedAt:       res.UpdatedAt,
	}, nil
}

func (r *CampaignRepo) UpdateCampaignStatus(ctx context.Context, id, storeID uuid.UUID, status string) error {
	params := generated.UpdateCampaignStatusParams{
		ID:       id,
		Status:   status,
		StoreID:  storeID,
	}
	return r.queries.UpdateCampaignStatus(ctx, params)
}
