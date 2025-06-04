package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	GeminiAPIKey string  `mapstructure:"gemini_api_key"`
	Model        string  `mapstructure:"model"`
	Temperature  float32 `mapstructure:"temperature"`
}

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".ted"), nil
}

func Load() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	viper.SetDefault("model", "gemini-2.0-flash")
	viper.SetDefault("temperature", 0.3)

	if err := os.MkdirAll(configPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := viper.SafeWriteConfig(); err != nil {
				return nil, fmt.Errorf("failed to write config file: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func Save(config *Config) error {
	viper.Set("gemini_api_key", config.GeminiAPIKey)
	viper.Set("model", config.Model)
	viper.Set("temperature", config.Temperature)

	return viper.WriteConfig()
}

func GetConfigPath() string {
	configPath, err := getConfigPath()
	if err != nil {
		return ""
	}
	return configPath
}
