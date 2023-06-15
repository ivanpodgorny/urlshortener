package security

import (
	"context"
)

// GRPCAuthenticator реализует методы для аутентификации и получения данных
// аутентифицированного пользователя из контекста GRPC-запроса.
type GRPCAuthenticator struct {
	tokenCreatorParser TokenCreatorParser
	userProvider       UserProvider[context.Context, context.Context]
}

// NewGRPCAuthenticator возвращает указатель на новый экземпляр GRPCAuthenticator.
func NewGRPCAuthenticator(cp TokenCreatorParser, p UserProvider[context.Context, context.Context]) *GRPCAuthenticator {
	return &GRPCAuthenticator{
		tokenCreatorParser: cp,
		userProvider:       p,
	}
}

// Authenticate получает идентификатор пользователя из токена, сохраненного в TokenStorage,
// и устанавливает его в UserProvider. Если токена не существует, генерирует новый идентификатор,
// формирует из него токен и сохраняет в TokenStorage.
func (a GRPCAuthenticator) Authenticate(ctx context.Context) context.Context {
	id, err := a.userProvider.Identifier(ctx)
	if err != nil {
		id = GenerateUUID()
	}
	ctx = a.userProvider.SetIdentifier(id, ctx)

	return ctx
}

// UserIdentifier возвращает идентификатор аутентифицированного пользователя из UserProvider.
func (a GRPCAuthenticator) UserIdentifier(ctx context.Context) (string, error) {
	return a.userProvider.Identifier(ctx)
}
