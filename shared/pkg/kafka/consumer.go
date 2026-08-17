package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

type ConsumerConfig struct {
	Brokers  []string
	Username string
	Password string
	CACert   string
	Topic    string
	GroupTopics []string
	GroupID  string
}

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(cfg ConsumerConfig) (*Consumer, error) {
	tlsCfg, err := buildTLSConfig(cfg.CACert)
	if err != nil {
		return nil, fmt.Errorf("build tls config: %w", err)
	}

	rc := kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		GroupID: cfg.GroupID,
		Dialer: &kafka.Dialer{
			SASLMechanism: plain.Mechanism{Username: cfg.Username, Password: cfg.Password},
			TLS:           tlsCfg,
		},
	}
	if len(cfg.GroupTopics) > 0 {
		rc.GroupTopics = cfg.GroupTopics
	} else {
		rc.Topic = cfg.Topic
	}

	return &Consumer{reader: kafka.NewReader(rc)}, nil
}

func (c *Consumer) Run(ctx context.Context, handler func(ctx context.Context, payload []byte) error) {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return
		}
		if err := handler(ctx, msg.Value); err != nil {
			continue
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}