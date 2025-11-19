package ds

import (
	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims представляет структуру JWT токена с пользовательскими данными
type JWTClaims struct {
	jwt.RegisteredClaims           // стандартные claims по RFC 7519
	UserID       int    `json:"user_id"`       // ID пользователя
	Username     string `json:"username"`      // имя пользователя
	IsModerator  bool   `json:"is_moderator"`  // флаг модератора
}

// JWTTokenResponse структура ответа с JWT токеном
type JWTTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"` // время жизни токена в секундах
}

