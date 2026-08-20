package config

import "os"

type Config struct {
	Port              string
	DatabaseURL       string
	PaystackSecretKey string
	RedisAddr         string
	PythonURL         string
}

func Load() *Config {
	return &Config{
		Port:              getEnv("PORT", "8080"),
		DatabaseURL:       getEnv("DATABASE_URL", "postgres://reset:resetpass@localhost:5432/reset?sslmode=disable"),
		PaystackSecretKey: getEnv("PAYSTACK_SECRET_KEY", ""),
		RedisAddr:         getEnv("REDIS_ADDR", "localhost:6379"),
		PythonURL:         getEnv("PYTHON_URL", "http://localhost:8001"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}