package middleware

import (
	"context"
	"errors"
	"strings"

	jwtManager "github.com/Yessenchik/bydz-fitness-app/services/auth-service/internal/infrastructure/jwt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func AuthInterceptor(
	jwt *jwtManager.JWTManager,
) grpc.UnaryServerInterceptor {

	publicMethods := map[string]bool{
		"/auth.v1.AuthService/Register":     true,
		"/auth.v1.AuthService/Login":        true,
		"/auth.v1.AuthService/RefreshToken": true,
	}

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		if publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("missing metadata")
		}

		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			return nil, errors.New("missing authorization header")
		}

		token := strings.TrimPrefix(
			authHeader[0],
			"Bearer ",
		)

		_, err := jwt.ValidateAccessToken(token)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}
