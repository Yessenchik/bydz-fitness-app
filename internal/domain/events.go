package domain

import (
	"time"

	"github.com/google/uuid"
)

// Имена NATS-топиков — единственный источник правды для всего приложения
const (
	EventUserRegistered  = "user.registered"
	EventUserDeleted     = "user.deleted"
	EventUserRoleChanged = "user.role_changed"
	EventEmailVerified   = "user.email_verified"
)

// UserRegisteredPayload публикуется при успешной регистрации.
// Сервис Membership слушает это событие и создаёт пустую карточку абонемента.
type UserRegisteredPayload struct {
	UserID     uuid.UUID `json:"user_id"`
	Email      string    `json:"email"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Role       Role      `json:"role"`
	OccurredAt time.Time `json:"occurred_at"`
}

// UserDeletedPayload публикуется при удалении пользователя.
type UserDeletedPayload struct {
	UserID     uuid.UUID `json:"user_id"`
	OccurredAt time.Time `json:"occurred_at"`
}

// UserRoleChangedPayload публикуется при смене роли.
type UserRoleChangedPayload struct {
	UserID     uuid.UUID `json:"user_id"`
	OldRole    Role      `json:"old_role"`
	NewRole    Role      `json:"new_role"`
	OccurredAt time.Time `json:"occurred_at"`
}

// EmailVerifiedPayload публикуется при подтверждении email.
type EmailVerifiedPayload struct {
	UserID     uuid.UUID `json:"user_id"`
	Email      string    `json:"email"`
	OccurredAt time.Time `json:"occurred_at"`
}
