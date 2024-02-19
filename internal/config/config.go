package config

import (
	"fmt"
	"github.com/dimitryshirokov/simple-app/internal/internal_error"
	"github.com/dimitryshirokov/simple-app/internal/logger"
	"os"
	"strconv"
)

type Config struct {
	DbUrl            string
	DbMinConnections int
	DbMaxConnections int
	QueryTimeout     int
	HttpPort         int
}

func NewConfig() (*Config, error) {
	dbUrl, err := getRequiredStringFromEnv("DB_URL")
	if err != nil {
		return nil, internal_error.NewError("can't get database url from env variables", err, nil)
	}
	return &Config{
		DbUrl:            dbUrl,
		DbMinConnections: getIntFromEnv("DB_MIN_CONNECTIONS", 1),
		DbMaxConnections: getIntFromEnv("DB_MAX_CONNECTIONS", 5),
		QueryTimeout:     getIntFromEnv("QUERY_TIMEOUT", 30),
		HttpPort:         getIntFromEnv("HTTP_PORT", 80),
	}, nil
}

func getRequiredStringFromEnv(env string) (string, error) {
	result := os.Getenv(env)
	if result == "" {
		return "", internal_error.NewError(fmt.Sprintf("env variable \"%s\" is empty", env), nil, map[string]interface{}{
			"env_variable": env,
		})
	}
	return result, nil
}

func getIntFromEnv(env string, defaultValue int) int {
	result, err := strconv.Atoi(os.Getenv(env))
	if err != nil {
		logger.LogWarning(fmt.Sprintf("can't get env variable \"%s\" value", env), map[string]interface{}{
			"env_var_value": os.Getenv(env),
			"default_value": defaultValue,
		}, err)
		result = defaultValue
	}
	return result
}
