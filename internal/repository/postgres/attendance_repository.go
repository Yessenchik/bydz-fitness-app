package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	"membership-service/internal/domain/entity"
)

type AttendancePostgresRepository struct {
	db *sql.DB
}

func NewAttendancePostgresRepository(db *sql.DB) *AttendancePostgresRepository {
	return &AttendancePostgresRepository{
		db: db,
	}
}

func (r *AttendancePostgresRepository) Create(
	ctx context.Context,
	attendance *entity.Attendance,
) error {
	query := `
		INSERT INTO attendances (
			id,
			user_id,
			check_in_time,
			date,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		attendance.ID,
		attendance.UserID,
		attendance.CheckInTime,
		attendance.Date,
		attendance.CreatedAt,
	)

	return err
}

func (r *AttendancePostgresRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*entity.Attendance, error) {
	query := `
		SELECT
			id,
			user_id,
			check_in_time,
			date,
			created_at
		FROM attendances
		WHERE id = $1
	`

	var attendance entity.Attendance

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&attendance.ID,
		&attendance.UserID,
		&attendance.CheckInTime,
		&attendance.Date,
		&attendance.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &attendance, nil
}

func (r *AttendancePostgresRepository) GetByUserID(
	ctx context.Context,
	userID uuid.UUID,
	limit int,
	offset int,
) ([]entity.Attendance, error) {
	query := `
		SELECT
			id,
			user_id,
			check_in_time,
			date,
			created_at
		FROM attendances
		WHERE user_id = $1
		ORDER BY check_in_time DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	attendances := make([]entity.Attendance, 0)

	for rows.Next() {
		var attendance entity.Attendance

		err := rows.Scan(
			&attendance.ID,
			&attendance.UserID,
			&attendance.CheckInTime,
			&attendance.Date,
			&attendance.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		attendances = append(attendances, attendance)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return attendances, nil
}

func (r *AttendancePostgresRepository) CountByUserAndPeriod(
	ctx context.Context,
	userID uuid.UUID,
	from time.Time,
	to time.Time,
) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM attendances
		WHERE user_id = $1
		  AND check_in_time >= $2
		  AND check_in_time < $3
	`

	var count int64

	err := r.db.QueryRowContext(ctx, query, userID, from, to).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *AttendancePostgresRepository) CountAllByPeriod(
	ctx context.Context,
	from time.Time,
	to time.Time,
) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM attendances
		WHERE check_in_time >= $1
		  AND check_in_time < $2
	`

	var count int64

	err := r.db.QueryRowContext(ctx, query, from, to).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *AttendancePostgresRepository) CountTotalByUser(
	ctx context.Context,
	userID uuid.UUID,
) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM attendances
		WHERE user_id = $1
	`

	var count int64

	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *AttendancePostgresRepository) CountTotal(
	ctx context.Context,
) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM attendances
	`

	var count int64

	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
