package auth

import (
	"fmt"
	"time"

	"lab1/internal/app/config"
	"lab1/internal/app/ds"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT генерирует JWT токен для пользователя
func GenerateJWT(user *ds.User, cfg config.JWT) (string, error) {
	now := time.Now()
	expiresAt := now.Add(cfg.ExpiresIn)

	claims := &ds.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "lab4-service",
			Subject:   fmt.Sprintf("%d", user.ID),
		},
		UserID:      user.ID,
		Username:    user.Username,
		IsModerator: user.IsModerator,
	}

	token := jwt.NewWithClaims(cfg.SigningMethod, claims)

	tokenString, err := token.SignedString([]byte(cfg.SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ParseJWT парсит и валидирует JWT токен
func ParseJWT(tokenString string, cfg config.JWT) (*ds.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if token.Method != cfg.SigningMethod {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.SecretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*ds.JWTClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// GetTokenExpiration возвращает оставшееся время жизни токена
func GetTokenExpiration(claims *ds.JWTClaims) time.Duration {
	if claims.ExpiresAt == nil {
		return 0
	}
	return time.Until(claims.ExpiresAt.Time)
}

