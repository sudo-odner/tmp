package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerAddr   string
	Timeout      time.Duration
	MaxRetries   int
	AssemblerTTL time.Duration
}

func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists && val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) (int, error) {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal, nil
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, fmt.Errorf("invalid int value for %s: %w", key, err)
	}
	return val, nil
}

func getEnvDuration(key string, defaultVal time.Duration) (time.Duration, error) {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal, nil
	}
	d, err := time.ParseDuration(valStr)
	if err != nil {
		return 0, fmt.Errorf("invalid duration value for %s: %w", key, err)
	}
	return d, nil
}

func MustLoad() *Config {
	var cfg Config
	var err error

	cfg.ServerAddr = getEnv("CLIENT_SERVER_ADDR", "0.0.0.0:5008")

	maxRetries, err := getEnvInt("CLIENT_RETRIES", 3)
	if err != nil {
		log.Fatal("ERROR: faled parce config: %w", err)
	}
	cfg.MaxRetries = maxRetries

	timeout, err := getEnvDuration("CLIENT_TIMEOUT", 2*time.Second)
	if err != nil {
		log.Fatal("ERROR: faled parce config: %w", err)
	}
	cfg.Timeout = timeout

	assemblerTTL, err := getEnvDuration("CLIENT_ASSEMBLER_TTL", 6*time.Second)
	if err != nil {
		log.Fatal("ERROR: faled parce config: %w", err)
	}
	cfg.AssemblerTTL = assemblerTTL

	return &cfg
}
