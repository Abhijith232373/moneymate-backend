package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moneymate-2026/moneymate-backend/services/notification/internal/domain"
	"github.com/moneymate-2026/moneymate-backend/services/notification/sqlc/generated"
)

type DeliveryRepo struct {
	q *generated.Queries
}

func NewDeliveryRepo(pool *pgxpool.Pool) *DeliveryRepo {
	return &DeliveryRepo{q: generated.New(pool)}
}

func (r *DeliveryRepo) Insert(ctx context.Context, l *domain.DeliveryLog) error {
	return r.q.InsertDeliveryLog(ctx, generated.InsertDeliveryLogParams{
		InboxID:           l.InboxID,
		DeviceTokenID:     uuidPtrToPgtype(l.DeviceTokenID),
		ProviderMessageID: strPtr(l.ProviderMessageID),
		Status:            l.Status,
		ErrorCode:         strPtr(l.ErrorCode),
	})
}

var _ domain.DeliveryRepository = (*DeliveryRepo)(nil)
