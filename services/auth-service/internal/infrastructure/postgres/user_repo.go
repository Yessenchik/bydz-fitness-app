package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/Yessenchik/bydz-fitness-app/user-auth-service/internal/domain"
)

type userRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewUserRepository(db *sqlx.DB, logger *zap.Logger) domain.UserRepository {
	return &userRepository{db: db, logger: logger}
}

// ─────────────────────────────────────────
//  Внутренняя модель БД (db-тег = колонка)
// ─────────────────────────────────────────

type userRow struct {
	ID           uuid.UUID `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	FirstName    string    `db:"first_name"`
	LastName     string    `db:"last_name"`
	Phone        string    `db:"phone"`
	Role         string    `db:"role"`
	IsVerified   bool      `db:"is_verified"`
	IsDeleted    bool      `db:"is_deleted"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func toEntity(r *userRow) *domain.User {
	return &domain.User{
		ID:           r.ID,
		Email:        r.Email,
		PasswordHash: r.PasswordHash,
		FirstName:    r.FirstName,
		LastName:     r.LastName,
		Phone:        r.Phone,
		Role:         domain.Role(r.Role),
		IsVerified:   r.IsVerified,
		IsDeleted:    r.IsDeleted,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

// ─────────────────────────────────────────
//  Create — внутри транзакции
// ─────────────────────────────────────────

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	// Используем транзакцию: вставка user + вставка audit-записи атомарны
	return r.withTx(ctx, func(tx *sqlx.Tx) error {
		const q = `
			INSERT INTO users
				(id, email, password_hash, first_name, last_name, phone,
				 role, is_verified, is_deleted, created_at, updated_at)
			VALUES
				(:id, :email, :password_hash, :first_name, :last_name, :phone,
				 :role, :is_verified, :is_deleted, :created_at, :updated_at)`

		row := &userRow{
			ID:           user.ID,
			Email:        user.Email,
			PasswordHash: user.PasswordHash,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			Phone:        user.Phone,
			Role:         string(user.Role),
			IsVerified:   user.IsVerified,
			IsDeleted:    false,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
		}

		if _, err := tx.NamedExecContext(ctx, q, row); err != nil {
			return fmt.Errorf("create user: %w", err)
		}

		// Audit log внутри той же транзакции
		const audit = `
			INSERT INTO user_audit_log (user_id, action, occurred_at)
			VALUES ($1, 'registered', $2)`
		if _, err := tx.ExecContext(ctx, audit, user.ID, time.Now().UTC()); err != nil {
			return fmt.Errorf("create user audit: %w", err)
		}

		return nil
	})
}

// ─────────────────────────────────────────
//  GetByID
// ─────────────────────────────────────────

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	const q = `
		SELECT id, email, password_hash, first_name, last_name, phone,
		       role, is_verified, is_deleted, created_at, updated_at
		FROM users
		WHERE id = $1 AND is_deleted = false`

	var row userRow
	if err := r.db.GetContext(ctx, &row, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return toEntity(&row), nil
}

// ─────────────────────────────────────────
//  GetByEmail
// ─────────────────────────────────────────

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	const q = `
		SELECT id, email, password_hash, first_name, last_name, phone,
		       role, is_verified, is_deleted, created_at, updated_at
		FROM users
		WHERE email = $1 AND is_deleted = false`

	var row userRow
	if err := r.db.GetContext(ctx, &row, q, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return toEntity(&row), nil
}

// ─────────────────────────────────────────
//  Update
// ─────────────────────────────────────────

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	const q = `
		UPDATE users SET
			first_name    = :first_name,
			last_name     = :last_name,
			phone         = :phone,
			role          = :role,
			is_verified   = :is_verified,
			password_hash = :password_hash,
			updated_at    = :updated_at
		WHERE id = :id AND is_deleted = false`

	row := &userRow{
		ID:           user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Phone:        user.Phone,
		Role:         string(user.Role),
		IsVerified:   user.IsVerified,
		PasswordHash: user.PasswordHash,
		UpdatedAt:    user.UpdatedAt,
	}

	result, err := r.db.NamedExecContext(ctx, q, row)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

// ─────────────────────────────────────────
//  SoftDelete — транзакция: пометить + audit
// ─────────────────────────────────────────

func (r *userRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.withTx(ctx, func(tx *sqlx.Tx) error {
		const q = `
			UPDATE users SET is_deleted = true, updated_at = $1
			WHERE id = $2 AND is_deleted = false`

		result, err := tx.ExecContext(ctx, q, time.Now().UTC(), id)
		if err != nil {
			return fmt.Errorf("soft delete user: %w", err)
		}
		rows, _ := result.RowsAffected()
		if rows == 0 {
			return domain.ErrUserNotFound
		}

		const audit = `
			INSERT INTO user_audit_log (user_id, action, occurred_at)
			VALUES ($1, 'deleted', $2)`
		if _, err := tx.ExecContext(ctx, audit, id, time.Now().UTC()); err != nil {
			return fmt.Errorf("soft delete audit: %w", err)
		}
		return nil
	})
}

// ─────────────────────────────────────────
//  List с пагинацией
// ─────────────────────────────────────────

func (r *userRepository) List(ctx context.Context, filter domain.ListUsersFilter) ([]*domain.User, int, error) {
	// Базовый запрос
	baseWhere := "WHERE is_deleted = false"
	args := []any{}
	argIdx := 1

	if filter.Role != nil {
		baseWhere += fmt.Sprintf(" AND role = $%d", argIdx)
		args = append(args, string(*filter.Role))
		argIdx++
	}

	// Считаем total
	var total int
	countQ := "SELECT COUNT(*) FROM users " + baseWhere
	if err := r.db.GetContext(ctx, &total, countQ, args...); err != nil {
		return nil, 0, fmt.Errorf("list users count: %w", err)
	}

	// Основная выборка с LIMIT/OFFSET
	offset := (filter.Page - 1) * filter.PageSize
	listQ := fmt.Sprintf(`
		SELECT id, email, password_hash, first_name, last_name, phone,
		       role, is_verified, is_deleted, created_at, updated_at
		FROM users %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, baseWhere, argIdx, argIdx+1)

	args = append(args, filter.PageSize, offset)

	var rows []userRow
	if err := r.db.SelectContext(ctx, &rows, listQ, args...); err != nil {
		return nil, 0, fmt.Errorf("list users select: %w", err)
	}

	users := make([]*domain.User, len(rows))
	for i, row := range rows {
		r2 := row
		users[i] = toEntity(&r2)
	}
	return users, total, nil
}

// ─────────────────────────────────────────
//  Password Reset Token
// ─────────────────────────────────────────

func (r *userRepository) SavePasswordResetToken(ctx context.Context, t *domain.PasswordResetToken) error {
	const q = `
		INSERT INTO password_reset_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE
			SET token = EXCLUDED.token, expires_at = EXCLUDED.expires_at`
	_, err := r.db.ExecContext(ctx, q, t.UserID, t.Token, t.ExpiresAt)
	return err
}

func (r *userRepository) GetPasswordResetToken(ctx context.Context, token string) (*domain.PasswordResetToken, error) {
	const q = `SELECT user_id, token, expires_at FROM password_reset_tokens WHERE token = $1`
	var t domain.PasswordResetToken
	if err := r.db.GetContext(ctx, &t, q, token); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTokenInvalid
		}
		return nil, err
	}
	return &t, nil
}

func (r *userRepository) DeletePasswordResetToken(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM password_reset_tokens WHERE token = $1`, token)
	return err
}

// ─────────────────────────────────────────
//  Email Verification Token
// ─────────────────────────────────────────

func (r *userRepository) SaveEmailVerificationToken(ctx context.Context, t *domain.EmailVerificationToken) error {
	const q = `
		INSERT INTO email_verification_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE
			SET token = EXCLUDED.token, expires_at = EXCLUDED.expires_at`
	_, err := r.db.ExecContext(ctx, q, t.UserID, t.Token, t.ExpiresAt)
	return err
}

func (r *userRepository) GetEmailVerificationToken(ctx context.Context, token string) (*domain.EmailVerificationToken, error) {
	const q = `SELECT user_id, token, expires_at FROM email_verification_tokens WHERE token = $1`
	var t domain.EmailVerificationToken
	if err := r.db.GetContext(ctx, &t, q, token); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTokenInvalid
		}
		return nil, err
	}
	return &t, nil
}

func (r *userRepository) DeleteEmailVerificationToken(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM email_verification_tokens WHERE token = $1`, token)
	return err
}

// ─────────────────────────────────────────
//  withTx — вспомогательная обёртка транзакции
// ─────────────────────────────────────────

func (r *userRepository) withTx(ctx context.Context, fn func(tx *sqlx.Tx) error) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			r.logger.Error("rollback failed", zap.Error(rbErr))
		}
		return err
	}
	return tx.Commit()
}
