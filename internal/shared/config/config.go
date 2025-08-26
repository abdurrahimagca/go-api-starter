package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Casbin   CasbinConfig
	JWT      JWTConfig
}

type DatabaseConfig struct {
	URL string
}

type ServerConfig struct {
	Port string
	Host string
	Env  string
}

type CasbinConfig struct {
	ModelPath  string
	PolicyPath string
}

type JWTConfig struct {
	Secret string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/go_api_starter?sslmode=disable"),
		},
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Host: getEnv("HOST", "localhost"),
			Env:  getEnv("ENV", "development"),
		},
		Casbin: CasbinConfig{
			ModelPath:  getEnv("CASBIN_MODEL_PATH", "./configs/casbin_model.conf"),
			PolicyPath: getEnv("CASBIN_POLICY_PATH", "./configs/casbin_policy.csv"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your-secret-key-here"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}