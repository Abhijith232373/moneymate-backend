package usecases

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"

	"github.com/moneymate-2026/moneymate-backend/services/notification/internal/domain"
)

type NotificationUsecase struct {
	inboxRepo      domain.InboxRepository
	preferenceRepo domain.PreferenceRepository
	deviceRepo     domain.DeviceTokenRepository
	deliveryRepo   domain.DeliveryRepository
	notifier       domain.Notifier
}

func NewNotificationUsecase(
	inbox domain.InboxRepository,
	preferences domain.PreferenceRepository,
	devices domain.DeviceTokenRepository,
	deliveries domain.DeliveryRepository,
	notifier domain.Notifier,
) *NotificationUsecase {
	return &NotificationUsecase{
		inboxRepo:      inbox,
		preferenceRepo: preferences,
		deviceRepo:     devices,
		deliveryRepo:   deliveries,
		notifier:       notifier,
	}
}

// HandleEvent is called once per Kafka message. It is idempotent: a replayed
// event hits the inbox unique constraint and becomes a no-op.
func (uc *NotificationUsecase) HandleEvent(ctx context.Context, env domain.EventEnvelope) error {
	recipientID, err := uuid.Parse(env.RecipientID)
	if err != nil {
		return err
	}

	tmpl, ok := templateFor(env)
	if !ok {
		log.Printf("notification: unknown event type %s — skipping", env.EventType)
		return nil
	}

	// 1. Persist the inbox row first (dedup guard + nothing is ever lost).
	inboxID, err := uc.inboxRepo.Insert(ctx, &domain.InboxMessage{
		RecipientType: env.RecipientType,
		RecipientID:   recipientID,
		Category:      tmpl.category,
		Title:         tmpl.title,
		Body:          tmpl.body,
		Data:          dataFor(env),
		EventID:       env.EventID,
	})
	if err != nil {
		return err
	}
	if inboxID == uuid.Nil {
		return nil // duplicate event, already handled
	}

	// 2. Respect per-category opt-out.
	enabled, err := uc.preferenceRepo.Get(ctx, env.RecipientType, recipientID, tmpl.category)
	if err != nil {
		return err
	}
	if !enabled {
		return uc.inboxRepo.MarkSent(ctx, inboxID)
	}

	// 3. Push to every active device this recipient owns.
	tokens, err := uc.deviceRepo.ListActiveByRecipient(ctx, env.RecipientType, recipientID)
	if err != nil {
		return err
	}
	for _, t := range tokens {
		uc.dispatch(ctx, inboxID, t, tmpl.toPushMessage())
	}
	return uc.inboxRepo.MarkSent(ctx, inboxID)
}

// dispatch sends one push and records it in the delivery log. Dead tokens are
// deactivated automatically.
func (uc *NotificationUsecase) dispatch(ctx context.Context, inboxID uuid.UUID, t *domain.DeviceToken, msg domain.PushMessage) {
	result, err := uc.notifier.Send(ctx, t.Token, msg)
	if err != nil {
		if result != nil && result.IsPermanent {
			_ = uc.deviceRepo.Deactivate(ctx, t.ID)
			_ = uc.deliveryRepo.Insert(ctx, &domain.DeliveryLog{
				InboxID: inboxID, DeviceTokenID: &t.ID, Status: "dropped", ErrorCode: result.ErrorCode,
			})
			return
		}
		_ = uc.deliveryRepo.Insert(ctx, &domain.DeliveryLog{
			InboxID: inboxID, DeviceTokenID: &t.ID, Status: "failed",
		})
		return
	}
	_ = uc.deliveryRepo.Insert(ctx, &domain.DeliveryLog{
		InboxID: inboxID, DeviceTokenID: &t.ID, Status: "sent", ProviderMessageID: result.ProviderMessageID,
	})
}

// --- template mapping (single source of truth for message copy) ---

type notificationTemplate struct {
	category  domain.Category
	title     string
	body      string
	channelID string
	data      map[string]string
}

func (t notificationTemplate) toPushMessage() domain.PushMessage {
	return domain.PushMessage{Title: t.title, Body: t.body, Data: t.data, ChannelID: t.channelID}
}

func templateFor(env domain.EventEnvelope) (notificationTemplate, bool) {
	fill := func(s string) string { return render(s, env.Payload) }

	switch env.EventType {
	case "moneymate.bills.due":
		return notificationTemplate{domain.CategoryBillDue,
			fill("{bill_name} bill due"), fill("₹{amount} due {due_date}"),
			"moneymate_alerts", map[string]string{"screen": "recurring"}}, true
	case "moneymate.debt.created":
		return notificationTemplate{domain.CategoryDebt,
			fill("{name} added a debt"), fill("Confirm ₹{amount} from {name}"),
			"moneymate_alerts", map[string]string{"screen": "debt", "id": env.PayloadValue("debt_id")}}, true
	case "moneymate.debt.confirmed":
		return notificationTemplate{domain.CategoryDebt,
			"Debt confirmed", fill("{name} confirmed ₹{amount}"),
			"moneymate_alerts", map[string]string{"screen": "debt", "id": env.PayloadValue("debt_id")}}, true
	case "moneymate.wallet.transfer.completed":
		return notificationTemplate{domain.CategoryTransfer,
			fill("Received ₹{amount}"), fill("from @{handle}"),
			"moneymate_alerts", map[string]string{"screen": "transaction", "id": env.PayloadValue("tx_id")}}, true
	case "moneymate.wallet.topup.completed":
		return notificationTemplate{domain.CategoryTransfer,
			fill("Wallet credited ₹{amount}"), "Your wallet balance has been updated",
			"moneymate_alerts", map[string]string{"screen": "wallet"}}, true
	case "moneymate.payment.merchant.completed":
		return notificationTemplate{domain.CategoryMerchant,
			"Payment received", fill("₹{amount} via QR"),
			"moneymate_alerts", map[string]string{"screen": "dashboard"}}, true
	case "moneymate.merchant.commission.credited":
		return notificationTemplate{domain.CategoryMerchant,
			fill("Commission credited ₹{amount}"), "Your payout balance was updated",
			"moneymate_alerts", map[string]string{"screen": "earnings"}}, true
	case "moneymate.merchant.verified":
		return notificationTemplate{domain.CategoryMerchant,
			"Business verified!", "Your QR code is now live",
			"moneymate_alerts", map[string]string{"screen": "qr"}}, true
	case "moneymate.campaign.created":
		return notificationTemplate{domain.CategoryCampaign,
			"Campaign is live", fill("{campaign_title}"),
			"moneymate_alerts", map[string]string{"screen": "campaigns", "id": env.PayloadValue("campaign_id")}}, true
	case "moneymate.offer.redeemed":
		return notificationTemplate{domain.CategoryOffer,
			"Offer redeemed", fill("{code} was used by a customer"),
			"moneymate_alerts", map[string]string{"screen": "campaigns", "id": env.PayloadValue("campaign_id")}}, true
	default:
		return notificationTemplate{}, false
	}
}

// dataFor builds the push data payload (no money values in push data).
func dataFor(env domain.EventEnvelope) map[string]any {
	data := make(map[string]any, len(env.Payload))
	for k, v := range env.Payload {
		switch k {
		case "amount", "due_date", "name", "handle", "bill_name", "campaign_title", "code":
			continue
		}
		data[k] = v
	}
	return data
}

// render replaces {key} placeholders from the event payload.
func render(s string, payload map[string]any) string {
	for k, v := range payload {
		s = strings.ReplaceAll(s, "{"+k+"}", fmtAny(v))
	}
	return s
}

func fmtAny(v any) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(fmt.Sprint(v), " ", ""), "\n", ""))
}
