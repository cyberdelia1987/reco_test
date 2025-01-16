package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Http            HttpConfig           `mapstructure:"http"`
	ShutdownTimeout time.Duration        `mapstructure:"shutdown_timeout"`
	Logging         LoggingConfig        `mapstructure:"logging"`
	CircuitBreaker  CircuitBreakerConfig `mapstructure:"circuit_breaker"`
	Asana           AsanaConfig          `mapstructure:"asana"`
}

type HttpConfig struct {
	Addr string `mapstructure:"addr"`
}

type LoggingConfig struct {
	Level         string   `mapstructure:"level"`
	Output        []string `mapstructure:"output"`
	LogStackTrace bool     `mapstructure:"log_stack_trace"`
}

type CircuitBreakerConfig struct {
	Name        string        `mapstructure:"name"`
	Timeout     time.Duration `mapstructure:"timeout"`
	MaxRequests uint32        `mapstructure:"max_requests"`
	MaxFailures uint32        `mapstructure:"max_failures"`
}

type AsanaConfig struct {
	BaseURL     string `mapstructure:"base_url"`
	AccessToken string `mapstructure:"access_token"`
}

func ReadConfig(configPath string) (Config, error) {
	viperConfig := viper.New()
	viperConfig.SetConfigFile(configPath)
	err := viperConfig.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	var config Config

	err = viperConfig.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
