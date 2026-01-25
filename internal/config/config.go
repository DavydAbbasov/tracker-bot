package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Telegram  Telegram
	PostreSQL PgConfig
}
type PgConfig struct {
	Host    string `env:"HOST_DB"`
	Port    uint16 `env:"PORT_DB"`
	Name    string `env:"NAME_DB"`
	User    string `env:"USER_DB"`
	Pass    string `env:"PASSWORD_DB"`
	SSLMode string `env:"SSL_MODE" env-default:"disable"`
}
type Telegram struct {
	TelegramToken    string `env:"TELEGRAM_TOKEN"`
	TelegramBotDebug bool   `env:"TELEGRAM_BOT_DEBUG"`
}

const (
	envConfigPath     = "CONFIG_PATH"
	defaultConfigPath = ".env"
)

func ParseConfig() (*Config, error) {
	var cfg Config

	path := os.Getenv(envConfigPath) //prod
	if path == "" {
		path = defaultConfigPath // local
	}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	return &cfg, nil
}

func (c Config) PostgresDSN() string {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.PostreSQL.Host, c.PostreSQL.Port, c.PostreSQL.User, c.PostreSQL.Pass, c.PostreSQL.Name, c.PostreSQL.SSLMode,
	)
	return dsn
}
