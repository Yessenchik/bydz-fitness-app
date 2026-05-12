package entity

import (
	"time"

	"github.com/google/uuid"
)

type Attendance struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	CheckInTime time.Time
	Date        time.Time
	CreatedAt   time.Time
}

func NewAttendance(userID uuid.UUID) *Attendance {
	now := time.Now().UTC()

	date := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0,
		0,
		0,
		0,
		time.UTC,
	)

	return &Attendance{
		ID:          uuid.New(),
		UserID:      userID,
		CheckInTime: now,
		Date:        date,
		CreatedAt:   now,
	}
}
