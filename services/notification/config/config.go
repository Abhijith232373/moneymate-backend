package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	sharedconfig "github.com/moneymate-2026/moneymate-backend/shared/config"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	HTTPAddr string `mapstructure:"http_addr"`
}

type FCMConfig struct {
	ProjectID       string `mapstructure:"project_id"`
	CredentialsPath string `mapstructure:"credentials_path"`
}

type Config struct {
	Env                   string
	Server                ServerConfig `mapstructure:"server"`
	Database              sharedconfig.DatabaseConfig
	Kafka                 sharedconfig.KafkaConfig
	JWT                   sharedconfig.JWTConfig
	FCM                   FCMConfig `mapstructure:"fcm"`
	InternalServiceSecret string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()
	yamlPath := os.Getenv("CONFIG_PATH")
	if yamlPath == "" {
		yamlPath = "./config/config.yaml"
	}

	v := viper.New()
	v.SetConfigFile(yamlPath)
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Printf("Warning: failed to read config file: %v\n", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	cfg.Database = sharedconfig.LoadDatabaseConfig(v, "notification")
	cfg.JWT = sharedconfig.LoadJWTConfig(v)
	cfg.Kafka = sharedconfig.LoadKafkaConfig(v) // reads KAFKA_* env vars + CA cert file
	cfg.FCM = FCMConfig{
		ProjectID:       sharedconfig.Get("FIREBASE_PROJECT_ID", ""),
		CredentialsPath: sharedconfig.Get("FIREBASE_CREDENTIALS_PATH", ""),
	}
	cfg.Env = sharedconfig.Get("ENVIRONMENT", "dev")
	cfg.InternalServiceSecret = sharedconfig.MustGet("INTERNAL_SERVICE_SECRET")
	return &cfg, nil
}