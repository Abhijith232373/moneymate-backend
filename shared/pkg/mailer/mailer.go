// package mailer

// import (
// 	"context"
// 	"crypto/rand"
// 	"encoding/hex"
// 	"fmt"
// 	"net/smtp"
// 	"strings"
// )

// type Config struct {
// 	Host        string
// 	Port        int
// 	Username    string
// 	Password    string
// 	FromAddress string
// 	FromName    string
// }

// type Client struct {
// 	cfg Config
// }

// func New(cfg Config) *Client {
// 	return &Client{cfg: cfg}
// }
// func (c *Client) Send(ctx context.Context, to, subject, htmlBody, textBody string) error {
//     if containsCRLF(to) || containsCRLF(subject) {
//         return fmt.Errorf("invalid header value: contains line break")
//     }
//     addr := fmt.Sprintf("%s:%d", c.cfg.Host, c.cfg.Port)
//     auth := smtp.PlainAuth("", c.cfg.Username, c.cfg.Password, c.cfg.Host)

//     msg, err := buildMIMEMessage(c.cfg.FromName, c.cfg.FromAddress, to, subject, htmlBody, textBody)
//     if err != nil {
//         return fmt.Errorf("build mime message: %w", err)
//     }

//     errCh := make(chan error, 1)
//     go func() {
//         errCh <- smtp.SendMail(addr, auth, c.cfg.FromAddress, []string{to}, msg)
//     }()

//     select {
//     case err := <-errCh:
//         if err != nil {
//             return fmt.Errorf("smtp send: %w", err)
//         }
//         return nil
//     case <-ctx.Done():
//         return ctx.Err()
//     }
// }

// func containsCRLF(s string) bool {
//     return strings.ContainsAny(s, "\r\n")
// }

// func buildMIMEMessage(fromName, fromAddr, to, subject, htmlBody, textBody string) ([]byte, error) {
//     boundary, err := randomBoundary()
//     if err != nil {
//         return nil, err
//     }

//     var b strings.Builder
//     fmt.Fprintf(&b, "From: %s <%s>\r\n", fromName, fromAddr)
//     fmt.Fprintf(&b, "To: %s\r\n", to)
//     fmt.Fprintf(&b, "Subject: %s\r\n", subject)
//     b.WriteString("MIME-Version: 1.0\r\n")
//     fmt.Fprintf(&b, "Content-Type: multipart/alternative; boundary=\"%s\"\r\n\r\n", boundary)

//     fmt.Fprintf(&b, "--%s\r\n", boundary)
//     b.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n")
//     b.WriteString(textBody)
//     b.WriteString("\r\n\r\n")

//     fmt.Fprintf(&b, "--%s\r\n", boundary)
//     b.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n")
//     b.WriteString(htmlBody)
//     b.WriteString("\r\n\r\n")

//     fmt.Fprintf(&b, "--%s--\r\n", boundary)

//     return []byte(b.String()), nil
// }

// func randomBoundary() (string, error) {
//     buf := make([]byte, 16)
//     if _, err := rand.Read(buf); err != nil {
//         return "", err
//     }
//     return hex.EncodeToString(buf), nil
// }

package mailer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Config struct {
	APIKey      string
	FromAddress string
	FromName    string
}

type Client struct {
	cfg Config
}

func New(cfg Config) *Client {
	return &Client{cfg: cfg}
}

type resendPayload struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Html    string   `json:"html"`
	Text    string   `json:"text"`
}
func (c *Client) Send(ctx context.Context, toEmail, subject, htmlBody, textBody string) error {
	payload := map[string]interface{}{
		"personalizations": []map[string]interface{}{
			{"to": []map[string]string{{"email": toEmail}}},
		},
		"from":    map[string]string{"email": c.cfg.FromAddress, "name": c.cfg.FromName},
		"subject": subject,
		"content": []map[string]string{
			{"type": "text/plain", "value": textBody},
			{"type": "text/html", "value": htmlBody},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal email payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.sendgrid.com/v3/mail/send", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("sendgrid request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("sendgrid returned status %d", resp.StatusCode)
	}
	return nil
}