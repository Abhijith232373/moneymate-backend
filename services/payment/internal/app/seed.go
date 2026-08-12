package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/moneymate-2026/moneymate-backend/services/payment/internal/domain"
	apperrors "github.com/moneymate-2026/moneymate-backend/shared/pkg/errors"
)

// seedExternalSettlementAccount ensures exactly one external_settlement
// account exists, creating it on first run. Idempotent: safe to call on
// every startup, across every replica — the unique partial index guarantees
// only one instance's INSERT wins if two replicas race, and the loser
// falls back to reading the winner's row.
func seedExternalSettlementAccount(ctx context.Context, accounts domain.AccountRepository) (uuid.UUID, error) {
	existing, err := accounts.GetExternalSettlementAccount(ctx)
	if err == nil {
		return existing.ID, nil
	}
	if !errors.Is(err, apperrors.ErrNotFound) {
		return uuid.Nil, fmt.Errorf("check external settlement account: %w", err)
	}

	created, err := accounts.CreateExternalSettlementAccount(ctx)
	if err != nil {
		if errors.Is(err, apperrors.ErrAlreadyExists) {
			// Lost a race against another replica — fetch the winner's row.
			existing, getErr := accounts.GetExternalSettlementAccount(ctx)
			if getErr != nil {
				return uuid.Nil, fmt.Errorf("fetch after race: %w", getErr)
			}
			return existing.ID, nil
		}
		return uuid.Nil, fmt.Errorf("create external settlement account: %w", err)
	}
	return created.ID, nil
}