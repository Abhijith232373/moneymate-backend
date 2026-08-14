package fcm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/moneymate-2026/moneymate-backend/services/notification/internal/domain"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const fcmSendScope = "https://www.googleapis.com/auth/firebase.messaging"

type Client struct {
	projectID string
	http      *http.Client
}

func New(ctx context.Context, projectID, credentialsPath string) (*Client, error) {
	data, err := os.ReadFile(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("read service account: %w", err)
	}
	creds, err := google.CredentialsFromJSON(ctx, data, fcmSendScope)
	if err != nil {
		return nil, fmt.Errorf("parse service account: %w", err)
	}
	return &Client{
		projectID: projectID,
		http:      oauth2.NewClient(ctx, creds.TokenSource),
	}, nil
}

type notificationBody struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type androidConfig struct {
	Notification struct {
		ChannelID string `json:"channel_id"`
		Priority  string `json:"priority"`
	} `json:"notification"`
}

type apnsConfig struct {
	Headers map[string]string `json:"headers"`
}

type fcmRequest struct {
	Message struct {
		Token        string            `json:"token"`
		Notification *notificationBody `json:"notification,omitempty"`
		Data         map[string]string `json:"data,omitempty"`
		Android      *androidConfig    `json:"android,omitempty"`
		APNS         *apnsConfig       `json:"apns,omitempty"`
	} `json:"message"`
}

func (c *Client) Send(ctx context.Context, token string, msg domain.PushMessage) (*domain.PushResult, error) {
	var req fcmRequest
	req.Message.Token = token
	req.Message.Notification = &notificationBody{Title: msg.Title, Body: msg.Body}
	req.Message.Data = msg.Data
	req.Message.Android = &androidConfig{}
	req.Message.Android.Notification.ChannelID = msg.ChannelID
	req.Message.Android.Notification.Priority = "HIGH"
	req.Message.APNS = &apnsConfig{Headers: map[string]string{"apns-priority": "10"}}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal fcm request: %w", err)
	}

	url := fmt.Sprintf("https://fcm.googleapis.com/v1/projects/%s/messages:send", c.projectID)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build fcm request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("fcm request: %w", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read fcm response: %w", err)
	}

	if resp.StatusCode == http.StatusOK {
		var ok struct {
			Name string `json:"name"`
		}
		_ = json.Unmarshal(respBody, &ok)
		return &domain.PushResult{ProviderMessageID: ok.Name}, nil
	}

	var errBody struct {
		Error struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		} `json:"error"`
	}
	_ = json.Unmarshal(respBody, &errBody)
	return &domain.PushResult{
		ErrorCode:   errBody.Error.Status,
		IsPermanent: isPermanent(resp.StatusCode, errBody.Error.Status),
	}, fmt.Errorf("fcm %s: %s", resp.Status, errBody.Error.Message)
}

// isPermanent reports whether the token must be deactivated (dead/unregistered)
// rather than retried.
func isPermanent(status int, grpcStatus string) bool {
	switch status {
	case http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound:
		return true
	}
	switch grpcStatus {
	case "UNREGISTERED", "INVALID_ARGUMENT", "NOT_FOUND":
		return true
	}
	return false // 429 / 5xx → transient, retry later
}
