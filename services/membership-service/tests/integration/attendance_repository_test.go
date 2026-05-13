package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"membership-service/internal/domain/entity"
	"membership-service/internal/repository/postgres"
)

func TestAttendanceRepository_CreateAndGetByUserID(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()

	repo := postgres.NewAttendancePostgresRepository(db)

	userID := uuid.New()

	attendance := entity.NewAttendance(userID)

	err := repo.Create(ctx, attendance)
	require.NoError(t, err)

	history, err := repo.GetByUserID(ctx, userID, 10, 0)
	require.NoError(t, err)

	require.Len(t, history, 1)
	assert.Equal(t, userID, history[0].UserID)
}

func TestAttendanceRepository_CountTotalByUser(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()

	repo := postgres.NewAttendancePostgresRepository(db)

	userID := uuid.New()

	firstAttendance := entity.NewAttendance(userID)
	secondAttendance := entity.NewAttendance(userID)

	err := repo.Create(ctx, firstAttendance)
	require.NoError(t, err)

	err = repo.Create(ctx, secondAttendance)
	require.NoError(t, err)

	count, err := repo.CountTotalByUser(ctx, userID)
	require.NoError(t, err)

	assert.Equal(t, int64(2), count)
}
