package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	MaxIPRequestsPerSecond    int
	MaxTokenRequestsPerSecond map[string]int
	RedisURL                  string
	IPLockTimeoutInSeconds    int
	TokenLockTimeoutInSeconds int
}

func NewConfig() *Config {
	godotenv.Load()
	return &Config{
		MaxIPRequestsPerSecond: getEnvAsInt("MAX_IP_REQUESTS", 10),
		MaxTokenRequestsPerSecond: map[string]int{
			os.Getenv("TOKEN_1"): getEnvAsInt("MAX_TOKEN1_REQUESTS", 20),
			os.Getenv("TOKEN_2"): getEnvAsInt("MAX_TOKEN2_REQUESTS", 30),
			os.Getenv("TOKEN_3"): getEnvAsInt("MAX_TOKEN3_REQUESTS", 50),
			os.Getenv("TOKEN_GOD"): getEnvAsInt("MAX_TOKEN_GOD_REQUESTS", 1000000),
		},
		RedisURL:                  os.Getenv("REDIS_URL"),
		IPLockTimeoutInSeconds:    getEnvAsInt("IP_LOCK_TIMEOUT", 10),
		TokenLockTimeoutInSeconds: getEnvAsInt("TOKEN_LOCK_TIMEOUT", 10),
	}
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
