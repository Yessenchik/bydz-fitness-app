package postgres

import (
	"fmt"
	"time"

	"github.com/Yessenchik/bydz-fitness-app/user-auth-service/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type jwtManager struct {
	accessSecret  string
	refreshSecret string
	accessTTL     time.Duration // 15 минут
	refreshTTL    time.Duration // 7 дней
}

func NewJWTManager(accessSecret, refreshSecret string) domain.TokenManager {
	return &jwtManager{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		accessTTL:     15 * time.Minute,
		refreshTTL:    7 * 24 * time.Hour,
	}
}

type jwtClaims struct {
	UserID  string `json:"uid"`
	Email   string `json:"email"`
	Role    string `json:"role"`
	TokenID string `json:"jti"`
	jwt.RegisteredClaims
}

func (m *jwtManager) GenerateTokenPair(user *domain.User) (*domain.Token, error) {
	now := time.Now().UTC()

	// Access token
	accessClaims := jwtClaims{
		UserID:  user.ID.String(),
		Email:   user.Email,
		Role:    string(user.Role),
		TokenID: uuid.New().String(),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTTL)),
			Issuer:    "user-auth-service",
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).
		SignedString([]byte(m.accessSecret))
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}

	// Refresh token (отдельный секрет — компрометация одного не ломает другой)
	refreshClaims := jwtClaims{
		UserID:  user.ID.String(),
		TokenID: uuid.New().String(),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.refreshTTL)),
			Issuer:    "user-auth-service",
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).
		SignedString([]byte(m.refreshSecret))
	if err != nil {
		return nil, fmt.Errorf("sign refresh token: %w", err)
	}

	return &domain.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(m.accessTTL.Seconds()),
	}, nil
}

func (m *jwtManager) ParseAccessToken(tokenString string) (*domain.TokenClaims, error) {
	return m.parse(tokenString, m.accessSecret)
}

func (m *jwtManager) ParseRefreshToken(tokenString string) (*domain.TokenClaims, error) {
	return m.parse(tokenString, m.refreshSecret)
}

func (m *jwtManager) parse(tokenString, secret string) (*domain.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{},
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secret), nil
		},
	)
	if err != nil {
		return nil, domain.ErrTokenInvalid
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return nil, domain.ErrTokenInvalid
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, domain.ErrTokenInvalid
	}

	return &domain.TokenClaims{
		UserID:  userID,
		Email:   claims.Email,
		Role:    domain.Role(claims.Role),
		TokenID: claims.TokenID,
	}, nil
}
