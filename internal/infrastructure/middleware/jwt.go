package middleware

import (
	"context"
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type JWTMiddleware struct {
	secret string
}

func NewJWTMiddleware(secret string) *JWTMiddleware {
	return &JWTMiddleware{
		secret: secret,
	}
}

func (m *JWTMiddleware) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if isPublicMethod(info.FullMethod) {
			return handler(ctx, req)
		}

		claims, err := m.extractClaims(ctx)
		if err != nil {
			return nil, err
		}

		userID, _ := claims["user_id"].(string)
		role, _ := claims["role"].(string)
		email, _ := claims["email"].(string)

		if userID == "" {
			return nil, errors.New("missing user_id in token")
		}

		if role == "" {
			return nil, errors.New("missing role in token")
		}

		ctx = context.WithValue(ctx, UserIDKey, userID)
		ctx = context.WithValue(ctx, RoleKey, role)
		ctx = context.WithValue(ctx, EmailKey, email)

		return handler(ctx, req)
	}
}

func (m *JWTMiddleware) extractClaims(ctx context.Context) (jwt.MapClaims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("missing metadata")
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return nil, errors.New("missing authorization header")
	}

	tokenString := strings.TrimSpace(values[0])
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	if tokenString == "" {
		return nil, errors.New("empty token")
	}

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}

			return []byte(m.secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func isPublicMethod(fullMethod string) bool {
	return true
	publicMethods := map[string]bool{
		"/membership.MembershipAccessService/CheckAccess": true,
	}

	return publicMethods[fullMethod]
}
