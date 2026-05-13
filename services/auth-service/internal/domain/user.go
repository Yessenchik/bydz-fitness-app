package domain

import "time"

type Role string

const (
	RoleClient  Role = "client"
	RoleTrainer Role = "trainer"
	RoleAdmin   Role = "admin"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	Phone        string
	Role         Role
	IsVerified   bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
