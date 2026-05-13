package domain

import "errors"

// Sentinel-ошибки предметной области.
// Usecase возвращает их; Delivery слой транслирует в gRPC-статусы.
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenExpired       = errors.New("token expired")
	ErrTokenInvalid       = errors.New("token invalid")
	ErrEmailNotVerified   = errors.New("email not verified")
	ErrPermissionDenied   = errors.New("permission denied")
	ErrWeakPassword       = errors.New("password too weak")
)
