package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"lab1/internal/app/config"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

const servicePrefix = "lab4_service." // префикс для всех ключей сервиса
const sessionPrefix = servicePrefix + "session." // префикс для ключей сессий

// Client представляет клиент Redis для работы с сессиями и blacklist
type Client struct {
	cfg    config.Redis
	client *redis.Client
}

// New создает новый Redis клиент
func New(ctx context.Context, cfg config.Redis) (*Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:    cfg.Password,
		DB:          cfg.DB,
		DialTimeout: cfg.DialTimeout,
		ReadTimeout: cfg.ReadTimeout,
	})

	// Проверяем соединение
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	logrus.Info("Redis connected successfully")

	return &Client{
		cfg:    cfg,
		client: redisClient,
	}, nil
}

// Close закрывает соединение с Redis
func (c *Client) Close() error {
	return c.client.Close()
}

// GetClient возвращает внутренний клиент Redis для прямого использования
func (c *Client) GetClient() *redis.Client {
	return c.client
}

// SaveSession сохраняет сессию пользователя в Redis (ключ - userID, значение - токен)
func (c *Client) SaveSession(ctx context.Context, token string, userID int, duration time.Duration) error {
	key := sessionPrefix + strconv.Itoa(userID)
	err := c.client.Set(ctx, key, token, duration).Err()
	if err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}
	logrus.Debugf("Session saved: userID %d -> token length: %d", userID, len(token))
	return nil
}

// GetSession получает ID пользователя по ID сессии
func (c *Client) GetSession(ctx context.Context, sessionID string) (int, error) {
	key := servicePrefix + "session." + sessionID
	result, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, fmt.Errorf("session not found")
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get session: %w", err)
	}

	userID, err := strconv.Atoi(result)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID in session: %w", err)
	}

	return userID, nil
}

// DeleteSession удаляет сессию из Redis
func (c *Client) DeleteSession(ctx context.Context, sessionID string) error {
	key := servicePrefix + "session." + sessionID
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	logrus.Debugf("Session deleted: %s", sessionID)
	return nil
}

// WriteJWTToBlacklist добавляет JWT токен в blacklist
func (c *Client) WriteJWTToBlacklist(ctx context.Context, jwtStr string, ttl time.Duration) error {
	key := servicePrefix + "jwt_blacklist." + jwtStr
	err := c.client.Set(ctx, key, true, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to add JWT to blacklist: %w", err)
	}
	logrus.Debugf("JWT added to blacklist: %s", jwtStr)
	return nil
}

// CheckJWTInBlacklist проверяет, находится ли JWT токен в blacklist
func (c *Client) CheckJWTInBlacklist(ctx context.Context, jwtStr string) (bool, error) {
	key := servicePrefix + "jwt_blacklist." + jwtStr
	err := c.client.Get(ctx, key).Err()
	if err == redis.Nil {
		// Токена нет в blacklist - это норма
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check JWT in blacklist: %w", err)
	}
	// Токен в blacklist
	return true, nil
}

// GetAllSessions возвращает все активные сессии (для демонстрации)
func (c *Client) GetAllSessions(ctx context.Context) (map[string]int, error) {
	pattern := servicePrefix + "session.*"
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get session keys: %w", err)
	}

	sessions := make(map[string]int)
	for _, key := range keys {
		value, err := c.client.Get(ctx, key).Result()
		if err != nil {
			continue
		}
		userID, err := strconv.Atoi(value)
		if err != nil {
			continue
		}
		sessions[key] = userID
	}

	return sessions, nil
}

