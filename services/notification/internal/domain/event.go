package domain

import (
	"fmt"
	"time"
)

type RecipientType string

const (
	RecipientUser     RecipientType = "user"
	RecipientMerchant RecipientType = "merchant"
)

type Category string

const (
	CategoryBillDue  Category = "bill_due"
	CategoryDebt     Category = "debt"
	CategoryTransfer Category = "transfer"
	CategoryMerchant Category = "merchant"
	CategoryCampaign Category = "campaign"
	CategoryOffer    Category = "offer"
	CategoryPromo    Category = "promo"
	CategorySystem   Category = "system"
)

// EventEnvelope is the contract every producer emits to the notification topics.
type EventEnvelope struct {
	EventID       string         `json:"event_id"`
	EventType     string         `json:"event_type"`
	OccurredAt    time.Time      `json:"occurred_at"`
	RecipientType RecipientType  `json:"recipient_type"`
	RecipientID   string         `json:"recipient_id"` // user_id or merchant_id
	Category      Category       `json:"category"`
	Payload       map[string]any `json:"payload"`
}

// PayloadValue returns the payload value for key as a string when present.
func (e EventEnvelope) PayloadValue(key string) string {
	v, ok := e.Payload[key]
	if !ok {
		return ""
	}
	return fmt.Sprint(v)
}