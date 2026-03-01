package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort     string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMode      string
	JWTSecret      string
	JWTExpiryHours int
	UploadPath      string
	MaxUploadSizeMB int
	BaseURL         string
	CORSOrigin      string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		ServerPort:      getEnv("SERVER_PORT", "8080"),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "postgres"),
		DBPassword:      getEnv("DB_PASSWORD", "postgres"),
		DBName:          getEnv("DB_NAME", "socialconnect"),
		DBSSLMode:       getEnv("DB_SSLMODE", "disable"),
		JWTSecret:       getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpiryHours:  getEnvAsInt("JWT_EXPIRY_HOURS", 72),
		UploadPath:      getEnv("UPLOAD_PATH", "./uploads"),
		MaxUploadSizeMB: getEnvAsInt("MAX_UPLOAD_SIZE_MB", 10),
		BaseURL:         getEnv("BASE_URL", "http://localhost:8080"),
		CORSOrigin:      getEnv("CORS_ORIGIN", "http://localhost:5173"),
	}
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
