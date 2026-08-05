package payment

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/razorpay/razorpay-go"
)

type RazorpayClient interface {
	CreateOrder(amount float64, currency string, receiptID string) (string, error)
	VerifySignature(orderID, paymentID, signature string) error
}

type razorpayClient struct {
	client *razorpay.Client
	secret string
}

func NewRazorpayClient(keyID, secret string) RazorpayClient {
	return &razorpayClient{
		client: razorpay.NewClient(keyID, secret),
		secret: secret,
	}
}

func (r *razorpayClient) CreateOrder(amount float64, currency string, receiptID string) (string, error) {
	// Amount should be in paise (smallest currency unit), so multiply by 100
	amountInPaise := int(amount * 100)

	data := map[string]interface{}{
		"amount":   amountInPaise,
		"currency": currency,
		"receipt":  receiptID,
	}

	body, err := r.client.Order.Create(data, nil)
	if err != nil {
		return "", err
	}

	orderID, ok := body["id"].(string)
	if !ok {
		return "", errors.New("failed to parse order id from razorpay response")
	}

	return orderID, nil
}

func (r *razorpayClient) VerifySignature(orderID, paymentID, signature string) error {
	data := orderID + "|" + paymentID

	h := hmac.New(sha256.New, []byte(r.secret))
	h.Write([]byte(data))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	if expectedSignature != signature {
		return errors.New("invalid signature")
	}

	return nil
}
