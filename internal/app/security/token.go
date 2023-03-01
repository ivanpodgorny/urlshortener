package security

import (
	"context"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
)

type userIDContextKey string

const (
	userIDCookie                  = "identity"
	userIDKey    userIDContextKey = "currentUserID"
)

var (
	ErrIncorrectHMACSignature = errors.New("incorrect hmac signature")
	ErrUserNotFound           = errors.New("user not found")
)

type TokenStorage[S, D any] interface {
	Get(source S) (string, bool)
	Set(token string, dest D)
}

type UserProvider[S, D any] interface {
	Identifier(source S) (string, error)
	SetIdentifier(id string, dest D) D
}

type TokenCreatorParser interface {
	Create(data string) string
	Parse(token string) (string, error)
}

type CookieTokenStorage struct {
	creatorParser TokenCreatorParser
}

func NewCookieTokenStorage(p TokenCreatorParser) *CookieTokenStorage {
	return &CookieTokenStorage{creatorParser: p}
}

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

func (s *CookieTokenStorage) Set(id string, w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:  userIDCookie,
		Value: s.creatorParser.Create(id),
		Path:  "/",
	})
}

type RequestContextUserProvider struct {
}

func (RequestContextUserProvider) Identifier(r *http.Request) (string, error) {
	val := r.Context().Value(userIDKey)
	if val == nil {
		return "", ErrUserNotFound
	}

	return val.(string), nil
}

func (RequestContextUserProvider) SetIdentifier(id string, r *http.Request) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), userIDKey, id))
}

type HMACTokenCreatorParser struct {
	key string
}

func NewHMACTokenCreatorParser(key string) *HMACTokenCreatorParser {
	return &HMACTokenCreatorParser{key: key}
}

func (cp *HMACTokenCreatorParser) Create(data string) string {
	return data + "/" + hex.EncodeToString(SignHMAC([]byte(data), cp.key))
}

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
