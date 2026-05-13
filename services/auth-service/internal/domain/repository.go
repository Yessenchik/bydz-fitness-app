package domain

import (
	"context"

	"github.com/google/uuid"
)

// ─────────────────────────────────────────
//  UserRepository — PostgreSQL
// ─────────────────────────────────────────

// UserRepository описывает все операции с хранилищем пользователей.
type UserRepository interface {
	// Создать пользователя (внутри транзакции, если tx != nil)
	Create(ctx context.Context, user *User) error

	// Получить по ID
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)

	// Получить по email
	GetByEmail(ctx context.Context, email string) (*User, error)

	// Обновить поля профиля
	Update(ctx context.Context, user *User) error

	// Мягкое удаление (is_deleted = true)
	SoftDelete(ctx context.Context, id uuid.UUID) error

	// Список пользователей с пагинацией и фильтром по роли
	List(ctx context.Context, filter ListUsersFilter) ([]*User, int, error)

	// Сохранить токен сброса пароля
	SavePasswordResetToken(ctx context.Context, t *PasswordResetToken) error

	// Найти токен сброса пароля
	GetPasswordResetToken(ctx context.Context, token string) (*PasswordResetToken, error)

	// Удалить использованный токен сброса
	DeletePasswordResetToken(ctx context.Context, token string) error

	// Сохранить токен верификации email
	SaveEmailVerificationToken(ctx context.Context, t *EmailVerificationToken) error

	// Найти токен верификации
	GetEmailVerificationToken(ctx context.Context, token string) (*EmailVerificationToken, error)

	// Удалить использованный токен верификации
	DeleteEmailVerificationToken(ctx context.Context, token string) error
}

// ListUsersFilter — параметры фильтрации и пагинации для ListUsers
type ListUsersFilter struct {
	Role     *Role // nil = все роли
	Page     int
	PageSize int
}

// ─────────────────────────────────────────
//  CacheRepository — Redis
// ─────────────────────────────────────────

// CacheRepository описывает операции с Redis-кэшем.
type CacheRepository interface {
	// Сохранить профиль пользователя в кэш (TTL задаётся реализацией)
	SetUserProfile(ctx context.Context, user *User) error

	// Получить профиль из кэша; (nil, nil) если нет записи
	GetUserProfile(ctx context.Context, id uuid.UUID) (*User, error)

	// Инвалидировать кэш профиля (после Update / Delete)
	InvalidateUserProfile(ctx context.Context, id uuid.UUID) error

	// Добавить access-токен в чёрный список (после Logout)
	// ttl = оставшееся время жизни токена
	BlacklistToken(ctx context.Context, tokenID string, ttlSeconds int64) error

	// Проверить, находится ли токен в чёрном списке
	IsTokenBlacklisted(ctx context.Context, tokenID string) (bool, error)

	// Сохранить refresh-токен (для ротации)
	SetRefreshToken(ctx context.Context, userID uuid.UUID, tokenID string, ttlSeconds int64) error

	// Получить refresh-токен по userID
	GetRefreshToken(ctx context.Context, userID uuid.UUID) (string, error)

	// Удалить refresh-токен (при Logout)
	DeleteRefreshToken(ctx context.Context, userID uuid.UUID) error
}

// ─────────────────────────────────────────
//  MessagePublisher — NATS
// ─────────────────────────────────────────

// MessagePublisher описывает публикацию событий в очередь сообщений.
// Реализация знает о NATS; Usecase работает только с этим интерфейсом.
type MessagePublisher interface {
	// Publish сериализует payload в JSON и отправляет в subject
	Publish(ctx context.Context, subject string, payload any) error
}

// ─────────────────────────────────────────
//  EmailService — SMTP
// ─────────────────────────────────────────

// EmailService описывает отправку писем.
// Реализация использует SMTP; Usecase работает только с этим интерфейсом.
type EmailService interface {
	// Отправить письмо с подтверждением email
	SendVerificationEmail(ctx context.Context, toEmail, firstName, token string) error

	// Отправить письмо со ссылкой для сброса пароля
	SendPasswordResetEmail(ctx context.Context, toEmail, firstName, token string) error
}

// ─────────────────────────────────────────
//  TokenManager — JWT
// ─────────────────────────────────────────

// TokenManager описывает генерацию и валидацию JWT.
// Позволяет подменить реализацию (например, на PASETO) без изменения Usecase.
type TokenManager interface {
	// Сгенерировать access + refresh токены для пользователя
	GenerateTokenPair(user *User) (*Token, error)

	// Разобрать и валидировать access-токен; вернуть claims
	ParseAccessToken(tokenString string) (*TokenClaims, error)

	// Разобрать и валидировать refresh-токен
	ParseRefreshToken(tokenString string) (*TokenClaims, error)
}

// TokenClaims — данные, закодированные в JWT
type TokenClaims struct {
	UserID  uuid.UUID
	Email   string
	Role    Role
	TokenID string // jti — уникальный ID токена (для блэклиста)
}
