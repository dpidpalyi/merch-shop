package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server struct {
		Port         string        `env:"SERVER_PORT"`
		ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT"`
		WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT"`
	}

	DB struct {
		Host     string `env:"DB_HOST"`
		Port     string `env:"DB_PORT"`
		User     string `env:"DB_USER"`
		Password string `env:"DB_PASSWORD"`
		Name     string `env:"DB_NAME"`
		SSLMode  string `env:"DB_SSL_MODE"`
	}

	//JWT struct {
	//	Secret        string
	//	TokenExpiry   time.Duration
	//	RefreshExpiry time.Duration
	//}
}

func New(path string) (*Config, error) {
	cfg := &Config{}

	cleanenv.ReadConfig(".env", cfg)

	err := cleanenv.ReadEnv(cfg)
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
