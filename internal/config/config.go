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
	Host        string `env-required:"true" yaml:"host"`
	Port        int    `env-required:"true" yaml:"port"`
	User        string `env-required:"true" yaml:"user"`
	Password    string `env-required:"true" env:"POSTGRES_PASSWORD"`
	Database    string `env-required:"true" yaml:"database"`
	SSLMode     string `env-required:"true" yaml:"sslmode"`
	DSNTemplate string
}

type Grpc struct {
	Host    string        `env-required:"true" yaml:"host"`
	Port    int           `env-required:"true" yaml:"port"`
	Timeout time.Duration `env-required:"true" yaml:"timeout"`
}

func MustLoad() *Config {
	err := godotenv.Load()
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
