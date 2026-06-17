package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort		string
	AppEnv		string
	JWTSecret	string

	DBHost 			string
	DBPort 			string
	DBUser 			string
	DBPassword 	string
	DBName 			string
	DBSSLMode 	string

	RedisHost 		string
	RedisPort 		string
	RedisPassword string
	RedisDB 			int

	CacheTTLDoctor 		int
	CacheTTLSchedule 	int
	CacheTTLQueue 		int
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		AppPort: 		getEnv("APP_PORT", "8080"),
		AppEnv: 		getEnv("APP_ENV", "development"),
		JWTSecret: 	mustGetEnv("JWT_SECRET"),

		DBHost: 		getEnv("DB_HOST", "localhost"),
		DBPort: 		getEnv("DB_PORT", "5433"),
		DBUser: 		mustGetEnv("DB_USER"),
		DBPassword: mustGetEnv("DB_PASSWORD"),
		DBName: 		mustGetEnv("DB_NAME"),
		DBSSLMode: 	getEnv("DB_SSLMODE", "diasble"),

		RedisHost: 			getEnv("REDIS_HOST", "localhost"),
		RedisPort: 			getEnv("REDIS_PORT", "6379"),
		RedisPassword: 	getEnv("REDIS_PASSWORD", ""),
		RedisDB: 				getEnvInt("REDIS_DB", 0),

		CacheTTLDoctor: getEnvInt("CACHE_TTL_DOCTOR", 300),
		CacheTTLSchedule: getEnvInt("CACHE_TTL_SCHEDULE", 600),
		CacheTTLQueue: getEnvInt("CACHE_TTL_QUEUE", 60),
	}, nil
}

func (c *Config) DBConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

func (c *Config) RedisAddr() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}

func getEnv(key, fallback string) string  {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return fallback
}

func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("required env var %q is not set", key))
	}
	return val
}

func getEnvInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return n
}