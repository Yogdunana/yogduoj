package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds the configuration for the judge service.
type Config struct {
	// GRPCAddr is the address the gRPC server listens on (e.g., ":50051").
	GRPCAddr string

	// PoolSize is the number of concurrent judging workers.
	PoolSize int

	// DBHost is the MySQL host address.
	DBHost string

	// DBPort is the MySQL port.
	DBPort int

	// DBUser is the MySQL username.
	DBUser string

	// DBPassword is the MySQL password.
	DBPassword string

	// DBName is the database name.
	DBName string
}

// Load reads configuration from environment variables and returns a Config.
// It provides sensible defaults for development.
func Load() (*Config, error) {
	cfg := &Config{
		GRPCAddr: getEnv("JUDGE_GRPC_ADDR", ":50051"),
		DBHost:   getEnv("DB_HOST", "localhost"),
		DBName:   getEnv("DB_NAME", "yogduoj"),
		DBUser:   getEnv("DB_USER", "root"),
	}

	poolSize, err := getEnvInt("JUDGE_POOL_SIZE", 4)
	if err != nil {
		return nil, fmt.Errorf("invalid JUDGE_POOL_SIZE: %w", err)
	}
	cfg.PoolSize = poolSize

	dbPort, err := getEnvInt("DB_PORT", 3306)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}
	cfg.DBPort = dbPort

	cfg.DBPassword = os.Getenv("DB_PASSWORD")

	return cfg, nil
}

// getEnv reads an environment variable or returns the fallback value.
func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

// getEnvInt reads an environment variable as an integer or returns the fallback value.
func getEnvInt(key string, fallback int) (int, error) {
	valStr := getEnv(key, "")
	if valStr == "" {
		return fallback, nil
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, err
	}
	return val, nil
}
