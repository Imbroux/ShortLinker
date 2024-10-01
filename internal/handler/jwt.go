package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
	"time"
)

// SIGNING_KEY - ключ для подписи JWT, загружается из переменной окружения
var SIGNING_KEY = []byte(os.Getenv("SIGNING_KEY"))

// Claims - структура для хранения информации о пользователе в JWT
type Claims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
}

// ErrInvalidAccessToken - ошибка, возвращаемая при недействительном токене
var ErrInvalidAccessToken = errors.New("invalid access token")

// MiddlewareJWT - middleware для проверки и обработки JWT токенов
func MiddlewareJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, "Unauthorized: Invalid token format", http.StatusUnauthorized)
			return
		}

		username, err := ParseToken(headerParts[1], SIGNING_KEY)
		if err != nil {
			status := http.StatusBadRequest
			if err == ErrInvalidAccessToken {
				status = http.StatusUnauthorized
			}
			http.Error(w, fmt.Sprintf("Unauthorized: %v", err), status)
			return
		}

		// Установка username в контекст запроса
		ctx := context.WithValue(r.Context(), "username", username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ParseToken - разбирает JWT токен и возвращает имя пользователя
func ParseToken(accessToken string, signKey []byte) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверка метода подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return signKey, nil
	})

	if err != nil {
		return "", err
	}

	// Проверка действительности токена
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.Username, nil
	}

	return "", ErrInvalidAccessToken
}

// CreateJWT - создает новый JWT токен для заданного пользователя
func CreateJWT(username string) (string, error) {
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Токен действителен 24 часа
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(SIGNING_KEY) // Подписание токена
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
