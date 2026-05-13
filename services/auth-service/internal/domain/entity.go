package domain

import (
	"time"

	"github.com/google/uuid"
)

// Role — ролевая модель пользователя
type Role string

const (
	RoleClient  Role = "client"
	RoleTrainer Role = "trainer"
	RoleAdmin   Role = "admin"
)

// User — корневая сущность предметной области.
// Никаких зависимостей от БД, gRPC или фреймворков.
type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	Phone        string
	Role         Role
	IsVerified   bool
	IsDeleted    bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Token — JWT-пара, возвращаемая при логине/рефреше
type Token struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64 // seconds
}

// PasswordResetToken — запись для сброса пароля
type PasswordResetToken struct {
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
}

// EmailVerificationToken — токен подтверждения почты
type EmailVerificationToken struct {
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
}
