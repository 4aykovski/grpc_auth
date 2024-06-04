package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env             string        `yaml:"env"`
	AccessTokenTtl  time.Duration `env-required:"true" yaml:"access_token_ttl"`
	RefreshTokenTtl time.Duration `env-required:"true" yaml:"refresh_token_ttl"`
	Postgres        Postgres      `env-required:"true" yaml:"postgres"`
	GRPC            Grpc          `env-required:"true" yaml:"grpc"`
}

type Postgres struct {
	Host        string `env-required:"true" env:"PG_HOST"`
	Port        int    `env-required:"true" env:"PG_PORT"`
	User        string `env-required:"true" env:"PG_USER"`
	Password    string `env-required:"true" env:"PG_PASS"`
	Database    string `env-required:"true" env:"PG_DBNAME"`
	SSLMode     string `env-required:"true" env:"PG_SSLMODE"`
	DSNTemplate string
}

type Grpc struct {
	Host    string        `env-required:"true" yaml:"host"`
	Port    int           `env-required:"true" yaml:"port"`
	Timeout time.Duration `env-required:"true" yaml:"timeout"`
}

// MustLoad loads config from .env and yaml file
//
// envPath is ".env" by default
func MustLoad(envPath string) *Config {
	if envPath == "" {
		envPath = ".env"
	}

	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("CONFIG_PATH %s does not exist", configPath)
	}

	cfg := &Config{}
	err = cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		log.Fatalf("cannot read config: %s", err.Error())
	}

	cfg.Postgres.DSNTemplate = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s", cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Database, cfg.Postgres.Password, cfg.Postgres.SSLMode)

	return cfg
}
