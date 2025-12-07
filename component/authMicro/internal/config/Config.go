package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPConfig HTTPConfig
	DB         DBConfig
	JWTConfig  JWTConfig
}

type HTTPConfig struct {
	ServerAddr string
}

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	AccessTokenExpiration  time.Duration
	RefreshTokenExpiration time.Duration
	Secret                 string
}

func NewConfig() (*Config, []error) {
	errors := make([]error, 0)
	return &Config{
		HTTPConfig: HTTPConfig{
			ServerAddr: getEnvIsRequired("ServerAddr", &errors),
		},
		DB: DBConfig{
			Host:     getEnvIsRequired("DBHost", &errors),
			Port:     getEnvIsRequired("DBPort", &errors),
			Username: getEnvIsRequired("DBUsername", &errors),
			Password: getEnvIsRequired("DBPassword", &errors),
			DBName:   getEnvIsRequired("DBName", &errors),
			SSLMode:  getEnvIsRequired("DBSSLMode", &errors),
		},
		JWTConfig: JWTConfig{
			AccessTokenExpiration:  time.Duration(getEnvIsRequiredAsInt("JWTAccessTokenExpiration", &errors)) * time.Millisecond,
			RefreshTokenExpiration: time.Duration(getEnvIsRequiredAsInt("JWTRefreshTokenExpiration", &errors)) * time.Millisecond,
			Secret:                 getEnvIsRequired("JWTSecret", &errors),
		},
	}, errors
}

func getEnv(key string) (string, error) {
	if value, exists := os.LookupEnv(key); exists {
		return value, nil
	}
	return "", fmt.Errorf("missing required environment variable: %s", key)
}

func getEnvIsRequired(key string, errors *[]error) string {
	result, e := getEnv(key)
	if e != nil {
		*errors = append(*errors, e)
		return ""
	}
	return result
}

func getEnvIsRequiredAsInt(key string, errors *[]error) int {
	resString, e := getEnv(key)
	if e != nil {
		*errors = append(*errors, e)
		return 0
	}
	result, err := strconv.Atoi(resString)
	if err != nil {
		*errors = append(*errors, err)
		return 0
	}
	return result
}
