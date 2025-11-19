package middleware

import (
	"net/http"
	"strings"

	"lab1/internal/app/auth"
	"lab1/internal/app/config"
	"lab1/internal/app/ds"
	"lab1/internal/app/redis"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	// AuthorizationHeader имя заголовка для авторизации
	AuthorizationHeader = "Authorization"
	// BearerPrefix префикс для Bearer токена
	BearerPrefix = "Bearer "
	// UserContextKey ключ для хранения пользователя в контексте
	UserContextKey = "user"
	// ClaimsContextKey ключ для хранения claims в контексте
	ClaimsContextKey = "claims"
)

// AuthMiddleware middleware для проверки JWT токена
type AuthMiddleware struct {
	cfg         config.JWT
	redisClient *redis.Client
}

// NewAuthMiddleware создает новый AuthMiddleware
func NewAuthMiddleware(cfg config.JWT, redisClient *redis.Client) *AuthMiddleware {
	return &AuthMiddleware{
		cfg:         cfg,
		redisClient: redisClient,
	}
}

// RequireAuth middleware, требующий авторизации
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем токен из заголовка
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, ds.ErrorResponse{
				Status:  "fail",
				Message: "authorization header is required",
			})
			c.Abort()
			return
		}

		// Проверяем префикс Bearer
		if !strings.HasPrefix(authHeader, BearerPrefix) {
			c.JSON(http.StatusUnauthorized, ds.ErrorResponse{
				Status:  "fail",
				Message: "invalid authorization header format",
			})
			c.Abort()
			return
		}

		// Извлекаем токен
		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)

		// Проверяем, не находится ли токен в blacklist
		isBlacklisted, err := m.redisClient.CheckJWTInBlacklist(c.Request.Context(), tokenString)
		if err != nil {
			logrus.Errorf("Failed to check JWT in blacklist: %v", err)
			c.JSON(http.StatusInternalServerError, ds.ErrorResponse{
				Status:  "error",
				Message: "internal server error",
			})
			c.Abort()
			return
		}

		if isBlacklisted {
			c.JSON(http.StatusUnauthorized, ds.ErrorResponse{
				Status:  "fail",
				Message: "token is revoked",
			})
			c.Abort()
			return
		}

		// Парсим и валидируем токен
		claims, err := auth.ParseJWT(tokenString, m.cfg)
		if err != nil {
			logrus.Errorf("Failed to parse JWT: %v", err)
			c.JSON(http.StatusUnauthorized, ds.ErrorResponse{
				Status:  "fail",
				Message: "invalid or expired token",
			})
			c.Abort()
			return
		}

		// Сохраняем claims в контексте для дальнейшего использования
		c.Set(ClaimsContextKey, claims)
		c.Set(UserContextKey, claims.UserID)

		c.Next()
	}
}

// RequireModerator middleware, требующий прав модератора
func (m *AuthMiddleware) RequireModerator() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Сначала проверяем авторизацию
		claimsValue, exists := c.Get(ClaimsContextKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, ds.ErrorResponse{
				Status:  "fail",
				Message: "authentication required",
			})
			c.Abort()
			return
		}

		claims, ok := claimsValue.(*ds.JWTClaims)
		if !ok {
			c.JSON(http.StatusInternalServerError, ds.ErrorResponse{
				Status:  "error",
				Message: "invalid claims format",
			})
			c.Abort()
			return
		}

		// Проверяем, является ли пользователь модератором
		if !claims.IsModerator {
			c.JSON(http.StatusForbidden, ds.ErrorResponse{
				Status:  "fail",
				Message: "moderator access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuth middleware для опциональной авторизации (не прерывает запрос, если токена нет)
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			// Нет токена - продолжаем без авторизации
			c.Next()
			return
		}

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)

		// Проверяем blacklist
		isBlacklisted, err := m.redisClient.CheckJWTInBlacklist(c.Request.Context(), tokenString)
		if err != nil || isBlacklisted {
			c.Next()
			return
		}

		// Парсим токен
		claims, err := auth.ParseJWT(tokenString, m.cfg)
		if err != nil {
			c.Next()
			return
		}

		// Сохраняем в контексте
		c.Set(ClaimsContextKey, claims)
		c.Set(UserContextKey, claims.UserID)

		c.Next()
	}
}

// GetUserID возвращает ID текущего пользователя из контекста
func GetUserID(c *gin.Context) (int, bool) {
	userID, exists := c.Get(UserContextKey)
	if !exists {
		return 0, false
	}
	id, ok := userID.(int)
	return id, ok
}

// GetClaims возвращает JWT claims из контекста
func GetClaims(c *gin.Context) (*ds.JWTClaims, bool) {
	claimsValue, exists := c.Get(ClaimsContextKey)
	if !exists {
		return nil, false
	}
	claims, ok := claimsValue.(*ds.JWTClaims)
	return claims, ok
}

// IsModerator проверяет, является ли текущий пользователь модератором
func IsModerator(c *gin.Context) bool {
	claims, ok := GetClaims(c)
	if !ok {
		return false
	}
	return claims.IsModerator
}

