package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

type Config struct {
	Brokers  []string
	Username string
	Password string
	CACert   string
}

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(cfg Config) (*Producer, error) {
	tlsCfg, err := buildTLSConfig(cfg.CACert)
	if err != nil {
		return nil, fmt.Errorf("build tls config: %w", err)
	}
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(cfg.Brokers...),
			Balancer: &kafka.LeastBytes{},
			Transport: &kafka.Transport{
				SASL: plain.Mechanism{Username: cfg.Username, Password: cfg.Password},
				TLS:  tlsCfg,
			},
		},
	}, nil
}

func (p *Producer) Publish(ctx context.Context, topic string, key, value []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{Topic: topic, Key: key, Value: value})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}