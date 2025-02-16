package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Server `yaml:"server"`
		DB     `yaml:"db"`
		JWT    `yaml:"jwt"`
	}

	Server struct {
		Port         string        `yaml:"port" env:"SERVER_PORT"`
		ReadTimeout  time.Duration `yaml:"read_timeout" env:"SERVER_READ_TIMEOUT"`
		WriteTimeout time.Duration `yaml:"write_timeout" env:"SERVER_WRITE_TIMEOUT"`
	}

	DB struct {
		Host     string `yaml:"host" env:"DB_HOST"`
		Port     string `yaml:"port" env:"DB_PORT"`
		User     string `yaml:"user" env:"DB_USER"`
		Password string `yaml:"password" env:"DB_PASSWORD"`
		Name     string `yaml:"name" env:"DB_NAME"`
		SSLMode  string `yaml:"ssl_mode" env:"DB_SSL_MODE"`
	}

	JWT struct {
		SecretKey   string        `yaml:"secret_key" env:"JWT_SECRET_KEY"`
		TokenExpiry time.Duration `yaml:"token_expiry" env:"JWT_TOKEN_EXPIRY"`
	}
)

func New(path string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("config.yml", cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DB.Host,
		c.DB.Port,
		c.DB.User,
		c.DB.Password,
		c.DB.Name,
		c.DB.SSLMode,
	)
}
