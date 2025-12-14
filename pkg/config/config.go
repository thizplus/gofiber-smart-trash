package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App     AppConfig
	DB      DatabaseConfig
	Storage StorageConfig
	AI      AIConfig
}

type AIConfig struct {
	ServiceURL string
	Timeout    int // in seconds
}

type AppConfig struct {
	Name string
	Port string
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type StorageConfig struct {
	Provider        string // r2, s3, gcs
	Bucket          string
	PublicURL       string
	PresignedExpiry int64 // in seconds

	// Cloudflare R2 / AWS S3 specific
	AccountID       string
	AccessKeyID     string
	SecretAccessKey string
	Region          string

	// GCS specific
	ProjectID       string
	CredentialsPath string
}

func LoadConfig() (*Config, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		// Don't fail if .env doesn't exist (e.g., in production with env vars)
		// return nil, err
	}

	presignedExpiry, _ := strconv.ParseInt(getEnv("PRESIGNED_URL_EXPIRY", "900"), 10, 64)
	aiTimeout, _ := strconv.Atoi(getEnv("AI_TIMEOUT", "30"))

	config := &Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "Smart Trash Picker API"),
			Port: getEnv("PORT", "3000"),
			Env:  getEnv("ENV", "development"),
		},
		AI: AIConfig{
			ServiceURL: getEnv("AI_SERVICE_URL", "http://localhost:8081"),
			Timeout:    aiTimeout,
		},
		DB: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "smartpicker"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Storage: StorageConfig{
			Provider:        getEnv("STORAGE_PROVIDER", "r2"),
			Bucket:          getEnv("R2_BUCKET", "smart-picker-bucket"),
			PublicURL:       getEnv("R2_PUBLIC_URL", ""),
			PresignedExpiry: presignedExpiry,

			// R2/S3
			AccountID:       getEnv("R2_ACCOUNT_ID", ""),
			AccessKeyID:     getEnv("R2_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("R2_SECRET_ACCESS_KEY", ""),
			Region:          getEnv("AWS_REGION", "auto"),

			// GCS (for future use)
			ProjectID:       getEnv("GCS_PROJECT_ID", ""),
			CredentialsPath: getEnv("GCS_CREDENTIALS_PATH", ""),
		},
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}