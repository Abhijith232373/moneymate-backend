package repo

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// uuidPtrToPgtype converts a nullable *uuid.UUID (domain model) into the
// pgtype.UUID that sqlc's generated code uses for nullable uuid columns.
func uuidPtrToPgtype(id *uuid.UUID) pgtype.UUID {
	if id == nil {
		return pgtype.UUID{}
	}
	return pgtype.UUID{Bytes: *id, Valid: true}
}

// timePtrToPgtype converts a nullable *time.Time (domain model) into the
// pgtype.Timestamptz that sqlc's generated code uses for nullable timestamps.
func timePtrToPgtype(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

// pgtypeToTimePtr converts a nullable pgtype.Timestamptz back to *time.Time.
func pgtypeToTimePtr(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}
