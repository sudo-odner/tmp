package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	MaxChunkSize int
	LossRate     float64
	Server       Server
	Cache        Cache
}

type Server struct {
	Address string
	Timeout time.Duration
}

type Cache struct {
	Capacity int
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

	maxChunkSize, err := getEnvInt("MAX_CHUNK_SIZE", 10)
	if err != nil {
		log.Fatalf("ERROR: failed to parse config: %v", err)
	}
	cfg.MaxChunkSize = maxChunkSize

	lossRate, err := getEnvInt("LOSS_RATE", 30)
	if err != nil {
		log.Fatalf("ERROR: failed to parse config: %v", err)
	}
	cfg.LossRate = float64(lossRate) / 100.0
	// Get server setting
	serverAddress := getEnv("SERVER_ADDRESS", "0.0.0.0:5008")
	serverTimeout, err := getEnvDuration("SERVER_TIMEOUT", 2*time.Second)
	if err != nil {
		log.Fatalf("ERROR: failed to parse config: %v", err)
	}
	cfg.Server = Server{
		Address: serverAddress,
		Timeout: serverTimeout,
	}

	// Get cache setting
	cacheCapacity, err := getEnvInt("CACHE_CAPACITY", 100)
	if err != nil {
		log.Fatalf("ERROR: failed to parse config: %v", err)
	}
	cfg.Cache = Cache{
		Capacity: cacheCapacity,
	}

	return &cfg
}
