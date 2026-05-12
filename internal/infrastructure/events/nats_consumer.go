package events

import (
	"encoding/json"
	"log/slog"

	"github.com/nats-io/nats.go"
)

type UserCreatedEvent struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

type UserDeletedEvent struct {
	UserID string `json:"user_id"`
}

type Consumer struct {
	conn   *nats.Conn
	logger *slog.Logger
}

func NewConsumer(conn *nats.Conn, logger *slog.Logger) *Consumer {
	return &Consumer{
		conn:   conn,
		logger: logger,
	}
}

func (c *Consumer) SubscribeUserCreated(handler func(event UserCreatedEvent) error) error {
	_, err := c.conn.Subscribe("user.created", func(msg *nats.Msg) {
		var event UserCreatedEvent

		if err := json.Unmarshal(msg.Data, &event); err != nil {
			c.logger.Error(
				"failed to unmarshal user.created event",
				"error", err,
			)
			return
		}

		if err := handler(event); err != nil {
			c.logger.Error(
				"failed to handle user.created event",
				"error", err,
				"user_id", event.UserID,
			)
			return
		}

		c.logger.Info(
			"user.created event handled",
			"user_id", event.UserID,
		)
	})

	return err
}

func (c *Consumer) SubscribeUserDeleted(handler func(event UserDeletedEvent) error) error {
	_, err := c.conn.Subscribe("user.deleted", func(msg *nats.Msg) {
		var event UserDeletedEvent

		if err := json.Unmarshal(msg.Data, &event); err != nil {
			c.logger.Error(
				"failed to unmarshal user.deleted event",
				"error", err,
			)
			return
		}

		if err := handler(event); err != nil {
			c.logger.Error(
				"failed to handle user.deleted event",
				"error", err,
				"user_id", event.UserID,
			)
			return
		}

		c.logger.Info(
			"user.deleted event handled",
			"user_id", event.UserID,
		)
	})

	return err
}
