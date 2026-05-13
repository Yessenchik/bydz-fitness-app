package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"membership-service/internal/domain/entity"
)

type AttendanceRepository interface {
	Create(ctx context.Context, attendance *entity.Attendance) error

	GetByID(ctx context.Context, id uuid.UUID) (*entity.Attendance, error)

	GetByUserID(ctx context.Context, userID uuid.UUID, limit int, offset int) ([]entity.Attendance, error)

	CountByUserAndPeriod(ctx context.Context, userID uuid.UUID, from time.Time, to time.Time) (int64, error)

	CountAllByPeriod(ctx context.Context, from time.Time, to time.Time) (int64, error)

	CountTotalByUser(ctx context.Context, userID uuid.UUID) (int64, error)

	CountTotal(ctx context.Context) (int64, error)
}
