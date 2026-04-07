package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	AppEnv         string
	Port           string
	JWTSecret      string
	MySQLDSN       string
	RedisAddr      string
	RedisPassword  string
	RedisDB        int
	FrontendOrigin string
}

func Load() Config {
	return Config{
		AppEnv:         getenv("APP_ENV", "development"),
		Port:           getenv("PORT", "8080"),
		JWTSecret:      getenv("JWT_SECRET", "change-me-in-production"),
		MySQLDSN:       getenv("MYSQL_DSN", "root:3217132@tcp(127.0.0.1:3306)/taoshop?parseTime=true&charset=utf8mb4&loc=Local"),
		RedisAddr:      getenv("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPassword:  getenv("REDIS_PASSWORD", ""),
		RedisDB:        mustAtoi(getenv("REDIS_DB", "0")),
		FrontendOrigin: getenv("FRONTEND_ORIGIN", "*"),
	}
}

func (c Config) Addr() string {
	return fmt.Sprintf(":%s", c.Port)
}

func getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}

func mustAtoi(value string) int {
	v, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return v
}
