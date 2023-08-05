package interceptor

import (
	"context"

	"google.golang.org/grpc"
)

// Authenticator интерфейс сервиса аутентификации пользователя.
type Authenticator interface {
	Authenticate(ctx context.Context) context.Context
}

// Authenticate возвращает interceptor для поверки токена пользователя.
func Authenticate(a Authenticator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = a.Authenticate(ctx)

		return handler(ctx, req)
	}
}
