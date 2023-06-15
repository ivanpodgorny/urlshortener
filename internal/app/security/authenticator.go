package security

import (
	"context"
	"net/http"
)

// Authenticator реализует методы для аутентификации и получения данных
// аутентифицированного пользователя.
type Authenticator struct {
	storage      TokenStorage[*http.Request, http.ResponseWriter]
	userProvider UserProvider[context.Context, *http.Request]
}

// NewAuthenticator возвращает указатель на новый экземпляр Authenticator.
func NewAuthenticator(s TokenStorage[*http.Request, http.ResponseWriter], p UserProvider[context.Context, *http.Request]) *Authenticator {
	return &Authenticator{
		storage:      s,
		userProvider: p,
	}
}

// Authenticate получает идентификатор пользователя из токена, сохраненного в TokenStorage,
// и устанавливает его в UserProvider. Если токена не существует, генерирует новый идентификатор,
// формирует из него токен и сохраняет в TokenStorage.
func (a Authenticator) Authenticate(w http.ResponseWriter, r *http.Request) *http.Request {
	id, ok := a.storage.Get(r)
	if !ok {
		id = a.generateID()
		a.storage.Set(id, w)
	}

	return a.userProvider.SetIdentifier(id, r)
}

// UserIdentifier возвращает идентификатор аутентифицированного пользователя из UserProvider.
func (a Authenticator) UserIdentifier(ctx context.Context) (string, error) {
	return a.userProvider.Identifier(ctx)
}

func (a Authenticator) generateID() string {
	return GenerateUUID()
}
