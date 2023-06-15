package security

import (
	"context"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type userIDContextKey string

const (
	userIDCookie                  = "identity"
	userIDKey    userIDContextKey = "currentUserID"
)

// ErrIncorrectHMACSignature ошибка проверки HMAC подписи.
var ErrIncorrectHMACSignature = errors.New("incorrect hmac signature")

// ErrUserNotFound ошибка получения данных аутентифицированного пользователя.
var ErrUserNotFound = errors.New("user not found")

// TokenStorage интерфейс сервиса для хранения аутентификационного токена.
type TokenStorage[S, D any] interface {
	Get(source S) (string, bool)
	Set(token string, dest D)
}

// UserProvider интерфейс сервиса для получения данных аутентифицированного пользователя.
type UserProvider[S, D any] interface {
	Identifier(source S) (string, error)
	SetIdentifier(id string, dest D) D
}

// TokenCreatorParser интерфейс сервиса создания и чтения аутентификационного токена.
type TokenCreatorParser interface {
	Create(data string) string
	Parse(token string) (string, error)
}

// CookieTokenStorage реализует методы для передачи и получения токена через Cookie.
type CookieTokenStorage struct {
	creatorParser TokenCreatorParser
}

// NewCookieTokenStorage возвращает указатель на новый экземпляр CookieTokenStorage.
func NewCookieTokenStorage(p TokenCreatorParser) *CookieTokenStorage {
	return &CookieTokenStorage{creatorParser: p}
}

// Get получает аутентификационный токен из Cookie.
func (s *CookieTokenStorage) Get(r *http.Request) (string, bool) {
	c, err := r.Cookie(userIDCookie)
	if err != nil {
		return "", false
	}

	id, err := s.creatorParser.Parse(c.Value)
	if err != nil {
		return "", false
	}

	return id, true
}

// Set устанавливает аутентификационный токен в Cookie.
func (s *CookieTokenStorage) Set(id string, w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:  userIDCookie,
		Value: s.creatorParser.Create(id),
		Path:  "/",
	})
}

// RequestContextUserProvider реализует методы для получения данных пользвателя из контекста запроса.
type RequestContextUserProvider struct {
}

// Identifier получает ID пользователя из контекста запроса.
func (RequestContextUserProvider) Identifier(ctx context.Context) (string, error) {
	val := ctx.Value(userIDKey)
	if val == nil {
		return "", ErrUserNotFound
	}

	return val.(string), nil
}

// SetIdentifier устанавливает ID пользователя в контекст запроса.
func (RequestContextUserProvider) SetIdentifier(id string, r *http.Request) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), userIDKey, id))
}

// GRPCContextUserProvider реализует методы для получения данных пользвателя из контекста GRPC-запроса.
type GRPCContextUserProvider struct {
}

// NewGRPCContextUserProvider возвращает указатель на новый экземпляр GRPCContextUserProvider.
func NewGRPCContextUserProvider() *GRPCContextUserProvider {
	return &GRPCContextUserProvider{}
}

// Identifier получает ID пользователя из контекста GRPC-запроса.
func (GRPCContextUserProvider) Identifier(ctx context.Context) (string, error) {
	val := ctx.Value(userIDKey)
	if val == nil {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok || len(md.Get(userIDCookie)) == 0 || md.Get(userIDCookie)[0] == "" {
			return "", ErrUserNotFound
		}

		return md.Get(userIDCookie)[0], nil
	}

	return val.(string), nil
}

// SetIdentifier устанавливает ID пользователя в контекст GRPC-запроса.
func (GRPCContextUserProvider) SetIdentifier(id string, ctx context.Context) context.Context {
	_ = grpc.SendHeader(ctx, metadata.Pairs(userIDCookie, id))

	return context.WithValue(ctx, userIDKey, id)
}

// HMACTokenCreatorParser реализует методы для создания и чтения аутентификационного токена.
type HMACTokenCreatorParser struct {
	key string
}

// NewHMACTokenCreatorParser возвращает указатель на новый экземпляр HMACTokenCreatorParser.
func NewHMACTokenCreatorParser(key string) *HMACTokenCreatorParser {
	return &HMACTokenCreatorParser{key: key}
}

// Create создает новый токен для аутентификации пользователя.
func (cp *HMACTokenCreatorParser) Create(data string) string {
	return data + "/" + hex.EncodeToString(SignHMAC([]byte(data), cp.key))
}

// Parse проверяет подлинность токена и получает из него значение ID.
func (cp *HMACTokenCreatorParser) Parse(token string) (string, error) {
	parts := strings.Split(token, "/")
	if len(parts) < 2 {
		return "", ErrIncorrectHMACSignature
	}

	data := parts[0]
	hmacSign, err := hex.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	if ValidateHMAC([]byte(data), hmacSign, cp.key) {
		return data, nil
	}

	return "", ErrIncorrectHMACSignature
}
