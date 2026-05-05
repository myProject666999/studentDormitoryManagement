package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Upload   UploadConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type JWTConfig struct {
	SecretKey string
	ExpireTime int64
}

type UploadConfig struct {
	MaxSize int64
	ImagePath string
	AttachmentPath string
}

var AppConfig *Config

func InitConfig() {
	AppConfig = &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "123456"),
			DBName:   getEnv("DB_NAME", "dormitory_management"),
		},
		JWT: JWTConfig{
			SecretKey:  getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			ExpireTime: getEnvAsInt64("JWT_EXPIRE", 86400),
		},
		Upload: UploadConfig{
			MaxSize:        getEnvAsInt64("UPLOAD_MAX_SIZE", 10*1024*1024),
			ImagePath:      getEnv("UPLOAD_IMAGE_PATH", "./uploads/images/"),
			AttachmentPath: getEnv("UPLOAD_ATTACHMENT_PATH", "./uploads/attachments/"),
		},
	}

	log.Println("Configuration initialized successfully")
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			return intValue
		}
	}
	return defaultValue
}
