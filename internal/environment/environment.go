package environment

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ResendEnvironment struct {
	Url string
	Key string
}

type TokenEnvironment struct {
	Secret                 string
	AccessTokenExpireTime  int
	RefreshTokenExpireTime int
	Issuer                 string
	Audience               string
}

type R2Environment struct {
	BucketName      string
	URL             string
	TokenValue      string
	AccessKeyID     string
	SecretAccessKey string
	AccountID       string
}

type Environment struct {
	APIKey      string
	Resend      ResendEnvironment
	DatabaseURL string
	Token       TokenEnvironment
	R2          R2Environment
	Port        string
}

func Load() (*Environment, error) {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	envFileName := fmt.Sprintf(".env.%s", env)
	if err := godotenv.Load(envFileName); err != nil {
		// Fallback to .env if specific env file not found
		if err := godotenv.Load(); err != nil {
			return nil, fmt.Errorf("error loading environment files: %w", err)
		}
	}

	accessTokenExpireTime := 3600 // default 1 hour
	if val := os.Getenv("ACCESS_TOKEN_EXPIRE_TIME"); val != "" {
		var err error
		accessTokenExpireTime, err = strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("error converting ACCESS_TOKEN_EXPIRE_TIME to int: %w", err)
		}
	}

	refreshTokenExpireTime := 604800 // default 7 days
	if val := os.Getenv("REFRESH_TOKEN_EXPIRE_TIME"); val != "" {
		var err error
		refreshTokenExpireTime, err = strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("error converting REFRESH_TOKEN_EXPIRE_TIME to int: %w", err)
		}
	}

	return &Environment{
		APIKey: os.Getenv("API_KEY"),
		Resend: ResendEnvironment{
			Url: os.Getenv("RESEND_URL"),
			Key: os.Getenv("RESEND_KEY"),
		},
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Token: TokenEnvironment{
			Secret:                 getEnvOrDefault("JWT_SECRET", "dev-secret-key-change-in-production"),
			AccessTokenExpireTime:  accessTokenExpireTime,
			RefreshTokenExpireTime: refreshTokenExpireTime,
			Issuer:                 getEnvOrDefault("ISSUER", "go-api-starter"),
			Audience:               getEnvOrDefault("AUDIENCE", "api-users"),
		},
		R2: R2Environment{
			BucketName:      os.Getenv("R2_BUCKET_NAME"),
			URL:             os.Getenv("R2_URL"),
			TokenValue:      os.Getenv("R2_TOKEN_VALUE"),
			AccessKeyID:     os.Getenv("R2_ACCESS_KEY_ID"),
			SecretAccessKey: os.Getenv("R2_SECRET_ACCESS_KEY"),
			AccountID:       os.Getenv("R2_ACCOUNT_ID"),
		},
		Port: getEnvOrDefault("PORT", "8080"),
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}