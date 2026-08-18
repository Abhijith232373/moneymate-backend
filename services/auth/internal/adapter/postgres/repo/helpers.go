package repo 
import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)



// timeToPgTimestamptz converts time.Time to pgtype.Timestamptz
func timeToPgTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}