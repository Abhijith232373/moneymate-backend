package domain

import "context"

type PushMessage struct {
	Title     string
	Body      string
	Data      map[string]string
	ChannelID string
}

type PushResult struct {
	ProviderMessageID string
	ErrorCode         string
	IsPermanent       bool // true = token dead (deactivate it)
}

type Notifier interface {
	Send(ctx context.Context, token string, msg PushMessage) (*PushResult, error)
}
