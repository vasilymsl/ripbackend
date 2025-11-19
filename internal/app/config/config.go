package config

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceHost string
	ServicePort int
	Database    Database
	Minio       Minio
	JWT         JWT
	Redis       Redis
}

type Database struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type Minio struct {
	Endpoint       string
	PublicEndpoint string `mapstructure:"public_endpoint"`
	AccessKey      string `mapstructure:"access_key"`
	SecretKey      string `mapstructure:"secret_key"`
	UseSSL         bool   `mapstructure:"use_ssl"`
	BucketName     string `mapstructure:"bucket_name"`
}

type JWT struct {
	SecretKey     string        `mapstructure:"secret_key"`
	ExpiresIn     time.Duration `mapstructure:"expires_in"` // время жизни токена
	SigningMethod jwt.SigningMethod
}

type Redis struct {
	Host        string        `mapstructure:"host"`
	Port        int           `mapstructure:"port"`
	Password    string        `mapstructure:"password"`
	DB          int           `mapstructure:"db"`
	DialTimeout time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout time.Duration `mapstructure:"read_timeout"`
}

func NewConfig() (*Config, error) {
	var err error

	configName := "config"
	_ = godotenv.Load()
	if os.Getenv("CONFIG_NAME") != "" {
		configName = os.Getenv("CONFIG_NAME")
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("toml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")
	viper.WatchConfig()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}           // создаем объект конфига
	err = viper.Unmarshal(cfg) // читаем информацию из файла,
	// конвертируем и затем кладем в нашу переменную cfg
	if err != nil {
		return nil, err
	}

	// Устанавливаем метод подписи JWT (по умолчанию HS256)
	cfg.JWT.SigningMethod = jwt.SigningMethodHS256

	// Устанавливаем значения по умолчанию, если не указаны
	if cfg.JWT.ExpiresIn == 0 {
		cfg.JWT.ExpiresIn = 24 * time.Hour // 24 часа по умолчанию
	}
	if cfg.Redis.DialTimeout == 0 {
		cfg.Redis.DialTimeout = 10 * time.Second
	}
	if cfg.Redis.ReadTimeout == 0 {
		cfg.Redis.ReadTimeout = 10 * time.Second
	}

	log.Info("config parsed")

	return cfg, nil
}
