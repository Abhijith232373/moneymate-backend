package paymentclient

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"time"
// )

// type Client struct {
// 	baseURL, secret string
// 	httpClient      *http.Client
// }

// func New(baseURL, secret string) *Client {
// 	return &Client{baseURL: baseURL, secret: secret, httpClient: &http.Client{Timeout: 10 * time.Second}}
// }

// func (c *Client) CreateWallet(ctx context.Context, userID, handle string) error {
// 	body, _ := json.Marshal(map[string]string{"user_id": userID, "handle": handle})
// 	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/internal/payment/wallets", bytes.NewReader(body))
// 	if err != nil {
// 		return fmt.Errorf("create request: %w", err)
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("X-Internal-Secret", c.secret)

// 	resp, err := c.httpClient.Do(req)
// 	if err != nil {
// 		return fmt.Errorf("payment-svc unreachable: %w", err)
// 	}
// 	defer resp.Body.Close()
// 	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
// 		return fmt.Errorf("payment-svc returned %d", resp.StatusCode)
// 	}
// 	return nil
// }