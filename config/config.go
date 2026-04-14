package config

import "os"

type Config struct {
	APIKey string
}

func NewConfig() *Config {
	config := &Config{}
	apiKey := os.Getenv("API_KEY")
	config.APIKey = apiKey
	return config
}
