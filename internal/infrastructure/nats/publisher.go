package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"

	"github.com/Yessenchik/bydz-fitness-app/user-auth-service/internal/domain"
)

type publisher struct {
	nc     *nats.Conn
	logger *zap.Logger
}

func NewPublisher(nc *nats.Conn, logger *zap.Logger) domain.MessagePublisher {
	return &publisher{nc: nc, logger: logger}
}

// Publish сериализует payload и публикует в NATS subject.
// Trace-контекст передаётся через заголовки сообщения (W3C TraceContext).
func (p *publisher) Publish(ctx context.Context, subject string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("nats publish marshal [%s]: %w", subject, err)
	}

	// Инжектируем trace-контекст в заголовки NATS-сообщения
	msg := &nats.Msg{
		Subject: subject,
		Data:    data,
		Header:  nats.Header{},
	}

	otel.GetTextMapPropagator().Inject(ctx, headerCarrier(msg.Header))

	if err := p.nc.PublishMsg(msg); err != nil {
		return fmt.Errorf("nats publish [%s]: %w", subject, err)
	}

	p.logger.Info("event published",
		zap.String("subject", subject),
		zap.Int("payload_bytes", len(data)),
	)
	return nil
}

// headerCarrier адаптирует nats.Header к интерфейсу propagation.TextMapCarrier
type headerCarrier map[string][]string

func (h headerCarrier) Get(key string) string {
	vals := h[key]
	if len(vals) == 0 {
		return ""
	}
	return vals[0]
}

func (h headerCarrier) Set(key, val string) {
	h[key] = []string{val}
}

func (h headerCarrier) Keys() []string {
	keys := make([]string, 0, len(h))
	for k := range h {
		keys = append(keys, k)
	}
	return keys
}
