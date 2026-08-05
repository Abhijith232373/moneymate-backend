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
		StoreID:        c.StoreID,
		Name:           c.Name,
		OfferType:      c.OfferType,
		RewardValue:    rewardValue,
		MinBillAmount:  minBillAmount,
		TargetAudience: c.TargetAudience,
		BannerUrl:      c.BannerURL,
		StartDate:      c.StartDate,
		EndDate:        c.EndDate,
		IsActive:       c.IsActive,
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
	c.OfferType = res.OfferType
	c.RewardValue = rv.Float64
	c.MinBillAmount = mb.Float64
	c.TargetAudience = res.TargetAudience
	c.BannerURL = res.BannerUrl
	c.StartDate = res.StartDate
	c.EndDate = res.EndDate
	c.IsActive = res.IsActive
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
			ID:             c.ID,
			StoreID:        c.StoreID,
			Name:           c.Name,
			OfferType:      c.OfferType,
			RewardValue:    rv.Float64,
			MinBillAmount:  mb.Float64,
			TargetAudience: c.TargetAudience,
			BannerURL:      c.BannerUrl,
			StartDate:      c.StartDate,
			EndDate:        c.EndDate,
			IsActive:       c.IsActive,
			CreatedAt:      c.CreatedAt,
			UpdatedAt:      c.UpdatedAt,
		})
	}
	return campaigns, nil
}

func (r *CampaignRepo) GetCampaignsByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*domain.Campaign, error) {
	res, err := r.queries.GetCampaignsByOwnerID(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	var campaigns []*domain.Campaign
	for _, c := range res {
		rv, _ := c.RewardValue.Float64Value()
		mb, _ := c.MinBillAmount.Float64Value()

		campaigns = append(campaigns, &domain.Campaign{
			ID:             c.ID,
			StoreID:        c.StoreID,
			Name:           c.Name,
			OfferType:      c.OfferType,
			RewardValue:    rv.Float64,
			MinBillAmount:  mb.Float64,
			TargetAudience: c.TargetAudience,
			BannerURL:      c.BannerUrl,
			StartDate:      c.StartDate,
			EndDate:        c.EndDate,
			IsActive:       c.IsActive,
			CreatedAt:      c.CreatedAt,
			UpdatedAt:      c.UpdatedAt,
		})
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
		ID:             res.ID,
		StoreID:        res.StoreID,
		Name:           res.Name,
		OfferType:      res.OfferType,
		RewardValue:    rv.Float64,
		MinBillAmount:  mb.Float64,
		TargetAudience: res.TargetAudience,
		BannerURL:      res.BannerUrl,
		StartDate:      res.StartDate,
		EndDate:        res.EndDate,
		IsActive:       res.IsActive,
		CreatedAt:      res.CreatedAt,
		UpdatedAt:      res.UpdatedAt,
	}, nil
}

func (r *CampaignRepo) UpdateCampaignStatus(ctx context.Context, id, storeID uuid.UUID, isActive bool) error {
	params := generated.UpdateCampaignStatusParams{
		ID:       id,
		IsActive: isActive,
		StoreID:  storeID,
	}
	return r.queries.UpdateCampaignStatus(ctx, params)
}
