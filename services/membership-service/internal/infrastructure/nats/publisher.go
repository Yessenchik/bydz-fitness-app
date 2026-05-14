package nats

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
)

type Publisher struct {
	conn *nats.Conn
	js   nats.JetStreamContext
}

func NewPublisher(url string) (*Publisher, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	js, err := conn.JetStream()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Publisher{
		conn: conn,
		js:   js,
	}, nil
}

func (p *Publisher) Publish(ctx context.Context, subject string, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = p.js.Publish(subject, payload)
	return err
}

func (p *Publisher) Close() {
	if p.conn != nil {
		p.conn.Close()
	}
}
