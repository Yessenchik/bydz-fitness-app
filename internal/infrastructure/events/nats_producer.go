package events

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
)

type Producer interface {
	Publish(ctx context.Context, subject string, payload any) error
}

type NATSProducer struct {
	conn *nats.Conn
}

func NewNATSProducer(conn *nats.Conn) *NATSProducer {
	return &NATSProducer{
		conn: conn,
	}
}

func (p *NATSProducer) Publish(ctx context.Context, subject string, payload any) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return p.conn.Publish(subject, data)
}
