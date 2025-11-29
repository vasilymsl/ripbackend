package service

import (
	"context"
	"fmt"
	"io"
	"lab1/internal/app/config"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

// MinioService - сервис для работы с MinIO
type MinioService struct {
	client         *minio.Client
	bucketName     string
	endpoint       string
	publicEndpoint string
}

// NewMinioService создает новый экземпляр MinIO сервиса
func NewMinioService(cfg config.Minio) (*MinioService, error) {
	// Инициализация клиента MinIO
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	// Используем публичный endpoint если он указан, иначе обычный
	publicEndpoint := cfg.PublicEndpoint
	if publicEndpoint == "" {
		publicEndpoint = cfg.Endpoint
	}

	service := &MinioService{
		client:         minioClient,
		bucketName:     cfg.BucketName,
		endpoint:       cfg.Endpoint,
		publicEndpoint: publicEndpoint,
	}

	// Создаем bucket если его нет
	if err := service.createBucketIfNotExists(); err != nil {
		return nil, fmt.Errorf("failed to create bucket: %w", err)
	}

	logrus.Info("MinIO service initialized successfully")
	return service, nil
}

// createBucketIfNotExists создает bucket если он не существует
func (s *MinioService) createBucketIfNotExists() error {
	ctx := context.Background()

	// Проверяем существование bucket
	exists, err := s.client.BucketExists(ctx, s.bucketName)
	if err != nil {
		return err
	}

	if !exists {
		// Создаем bucket
		err = s.client.MakeBucket(ctx, s.bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}

		// Устанавливаем политику публичного чтения для bucket
		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {"AWS": ["*"]},
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::%s/*"]
				}
			]
		}`, s.bucketName)

		err = s.client.SetBucketPolicy(ctx, s.bucketName, policy)
		if err != nil {
			logrus.Warnf("Failed to set bucket policy: %v", err)
		}

		logrus.Infof("Bucket '%s' created successfully", s.bucketName)
	}

	return nil
}

// UploadFile загружает файл в MinIO и возвращает имя файла
func (s *MinioService) UploadFile(file *multipart.FileHeader, folder string) (string, error) {
	ctx := context.Background()

	// Открываем файл
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Генерируем уникальное имя файла
	ext := filepath.Ext(file.Filename)
	fileName := generateFileName(ext)
	objectName := fmt.Sprintf("%s/%s", folder, fileName)

	// Определяем content type
	contentType := getContentType(ext)

	// Загружаем файл
	_, err = s.client.PutObject(ctx, s.bucketName, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	logrus.Infof("File uploaded successfully: %s", objectName)
	return fileName, nil
}

// DeleteFile удаляет файл из MinIO
func (s *MinioService) DeleteFile(folder, fileName string) error {
	if fileName == "" {
		return nil // Нечего удалять
	}

	ctx := context.Background()
	objectName := fmt.Sprintf("%s/%s", folder, fileName)

	err := s.client.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	logrus.Infof("File deleted successfully: %s", objectName)
	return nil
}

// GetFileURL возвращает публичный URL файла
func (s *MinioService) GetFileURL(folder, fileName string) string {
	if fileName == "" {
		return ""
	}
	
	// Если fileName уже полный URL (начинается с http:// или https://), возвращаем как есть
	if len(fileName) >= 7 && fileName[:7] == "http://" {
		return fileName
	}
	if len(fileName) >= 8 && fileName[:8] == "https://" {
		return fileName
	}
	
	objectName := fileName
	if folder != "" {
		objectName = fmt.Sprintf("%s/%s", folder, fileName)
	}
	return fmt.Sprintf("http://%s/%s/%s", s.publicEndpoint, s.bucketName, objectName)
}

// GetPresignedURL возвращает временную предподписанную ссылку на файл
func (s *MinioService) GetPresignedURL(folder, fileName string, expiry time.Duration) (string, error) {
	if fileName == "" {
		return "", nil
	}

	ctx := context.Background()
	objectName := fmt.Sprintf("%s/%s", folder, fileName)

	presignedURL, err := s.client.PresignedGetObject(ctx, s.bucketName, objectName, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedURL.String(), nil
}

// DownloadFile скачивает файл из MinIO
func (s *MinioService) DownloadFile(folder, fileName string) (io.ReadCloser, error) {
	ctx := context.Background()
	objectName := fmt.Sprintf("%s/%s", folder, fileName)

	object, err := s.client.GetObject(ctx, s.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	return object, nil
}

// generateFileName генерирует уникальное имя файла на латинице
func generateFileName(ext string) string {
	// Генерируем UUID и берем первые 12 символов
	id := uuid.New().String()[:12]
	timestamp := time.Now().Format("20060102")
	fileName := fmt.Sprintf("credit_%s_%s%s", timestamp, id, ext)
	return transliterate(fileName)
}

// transliterate транслитерирует имя файла в латиницу
func transliterate(input string) string {
	// Простая транслитерация (можно расширить)
	replacer := strings.NewReplacer(
		" ", "_",
		"а", "a", "б", "b", "в", "v", "г", "g", "д", "d",
		"е", "e", "ё", "yo", "ж", "zh", "з", "z", "и", "i",
		"й", "y", "к", "k", "л", "l", "м", "m", "н", "n",
		"о", "o", "п", "p", "р", "r", "с", "s", "т", "t",
		"у", "u", "ф", "f", "х", "h", "ц", "ts", "ч", "ch",
		"ш", "sh", "щ", "sch", "ъ", "", "ы", "y", "ь", "",
		"э", "e", "ю", "yu", "я", "ya",
		"А", "A", "Б", "B", "В", "V", "Г", "G", "Д", "D",
		"Е", "E", "Ё", "Yo", "Ж", "Zh", "З", "Z", "И", "I",
		"Й", "Y", "К", "K", "Л", "L", "М", "M", "Н", "N",
		"О", "O", "П", "P", "Р", "R", "С", "S", "Т", "T",
		"У", "U", "Ф", "F", "Х", "H", "Ц", "Ts", "Ч", "Ch",
		"Ш", "Sh", "Щ", "Sch", "Ъ", "", "Ы", "Y", "Ь", "",
		"Э", "E", "Ю", "Yu", "Я", "Ya",
	)
	return replacer.Replace(input)
}

// getContentType определяет content type по расширению файла
func getContentType(ext string) string {
	ext = strings.ToLower(ext)
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	default:
		return "application/octet-stream"
	}
}
