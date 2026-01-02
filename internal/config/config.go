package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Rate limiting by IP
	RateLimitIP      int
	IPBlockDuration  int

	// Rate limiting by Token
	RateLimitToken       int
	TokenBlockDuration   int

	// Redis configuration
	RedisHost string
	RedisPort int
	RedisDB   int

	// Server configuration
	ServerPort int
}

func LoadConfig() *Config {
	// Load .env file if it exists
	_ = godotenv.Load()

	return &Config{
		RateLimitIP:        getEnvAsInt("RATE_LIMIT_IP", 5),
		IPBlockDuration:    getEnvAsInt("IP_BLOCK_DURATION", 300),
		RateLimitToken:     getEnvAsInt("RATE_LIMIT_TOKEN", 10),
		TokenBlockDuration: getEnvAsInt("TOKEN_BLOCK_DURATION", 600),
		RedisHost:          getEnv("REDIS_HOST", "localhost"),
		RedisPort:          getEnvAsInt("REDIS_PORT", 6379),
		RedisDB:            getEnvAsInt("REDIS_DB", 0),
		ServerPort:         getEnvAsInt("SERVER_PORT", 8080),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}
