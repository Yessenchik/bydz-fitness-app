package nats

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
)

type Publisher struct {
	conn *nats.Conn
}

func NewPublisher(url string) (*Publisher, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return &Publisher{conn: conn}, nil
}

func (p *Publisher) Publish(subject string, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return p.conn.Publish(subject, payload)
}

func (p *Publisher) Close() {
	if p.conn != nil {
		p.conn.Close()
	}
}
