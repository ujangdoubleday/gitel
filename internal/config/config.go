package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// config holds all application configurations.
type Config struct {
	Server  ServerConfig
	Webhook WebhookConfig
	LLM     LLMConfig
}

// serverConfig holds HTTP server configuration.
type ServerConfig struct {
	Port string
}

// webhookConfig holds GitHub webhook configuration.
type WebhookConfig struct {
	Secret string
}

// llmConfig holds LLM provider configuration.
type LLMConfig struct {
	APIKey  string
	Model   string
	BaseURL string
	Timeout time.Duration
}

// load reads configuration from environment variables.
func Load() (*Config, error) {
	timeoutSec, _ := strconv.Atoi(getEnv("LLM_TIMEOUT_SECONDS", "30"))

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Webhook: WebhookConfig{
			Secret: getEnv("WEBHOOK_SECRET", ""),
		},
		LLM: LLMConfig{
			APIKey:  getEnv("LLM_API_KEY", ""),
			Model:   getEnv("LLM_MODEL", "gpt-4o-mini"),
			BaseURL: getEnv("LLM_BASE_URL", "https://api.openai.com/v1"),
			Timeout: time.Duration(timeoutSec) * time.Second,
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate checks if all required configurations are present.
func (c *Config) Validate() error {
	if c.Webhook.Secret == "" {
		return fmt.Errorf("WEBHOOK_SECRET is required")
	}
	if c.LLM.APIKey == "" {
		return fmt.Errorf("LLM_API_KEY is required")
	}
	if c.LLM.Model == "" {
		return fmt.Errorf("LLM_MODEL is required")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
