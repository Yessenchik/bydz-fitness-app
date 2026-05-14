package domain

import "time"

type Plan struct {
	ID           string
	Name         string
	DurationDays int32
	PriceKZT     int64
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
