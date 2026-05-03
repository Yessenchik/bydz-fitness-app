package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New создаёт production-ready структурированный логгер.
// В dev-режиме (isDev=true) выводит в консоль с цветами.
func New(isDev bool) (*zap.Logger, error) {
	var cfg zap.Config
	if isDev {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		cfg = zap.NewProductionConfig()
		// JSON-формат, читаемый Grafana Loki
		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}
	return cfg.Build()
}
