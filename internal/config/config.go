package config

import (
	"fmt"
	"os"
)

// config holds all application configurations.
type Config struct {
	Server  ServerConfig
	Webhook WebhookConfig
}

// serverConfig holds HTTP server configuration.
type ServerConfig struct {
	Port string
}

// webhookConfig holds GitHub webhook configuration.
type WebhookConfig struct {
	Secret string
}

// load reads configuration from environment variables.
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Webhook: WebhookConfig{
			Secret: getEnv("WEBHOOK_SECRET", ""),
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
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
