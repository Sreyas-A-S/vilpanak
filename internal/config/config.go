package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	Env             string
	AppURL          string
	SessionSecret   string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	DBSSLMode       string
	UploadDir       string
	MaxUploadSizeMB int
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Note: .env file not found, using system environment variables")
	}

	maxUpload, _ := strconv.Atoi(getEnv("MAX_UPLOAD_SIZE_MB", "10"))

	return &Config{
		Port:            getEnv("PORT", "8080"),
		Env:             getEnv("ENV", "development"),
		AppURL:          getEnv("APP_URL", "http://localhost:8080"),
		SessionSecret:   getEnv("SESSION_SECRET", "default-vilpanak-secret-32-bytes!"),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "postgres"),
		DBPassword:      getEnv("DB_PASSWORD", "postgres"),
		DBName:          getEnv("DB_NAME", "vilpanak_db"),
		DBSSLMode:       getEnv("DB_SSLMODE", "disable"),
		UploadDir:       getEnv("UPLOAD_DIR", "./static/uploads"),
		MaxUploadSizeMB: maxUpload,
	}
}

func (c *Config) DBConnectionString() string {
	return "postgres://" + c.DBUser + ":" + c.DBPassword + "@" + c.DBHost + ":" + c.DBPort + "/" + c.DBName + "?sslmode=" + c.DBSSLMode
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
