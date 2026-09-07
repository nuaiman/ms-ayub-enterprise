package config

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// go run ./cmd/app/ --config=config/config.env

type Config struct {
	Addr   string `env:"ADDR" env-default:":8090"`
	JWTKey string `env:"JWT_KEY"`

	RateLimit    int           `env:"RATE_LIMIT" env-default:"100"`
	RateInterval time.Duration `env:"RATE_INTERVAL" env-default:"1m"`

	DBPath     string `env:"DB_PATH" env-default:"./data"`
	DBName     string `env:"DB_NAME" env-default:"app.db"`
	SchemaPath string `env:"SCHEMA_PATH" env-default:"internal/db/schema/schema.sql"` // FIXED

	BootstrapAdmin     bool   `env:"BOOTSTRAP_ADMIN" env-default:"true"`
	BootstrapUsername  string `env:"BOOTSTRAP_ADMIN_USERNAME"`
	BootstrapPassword  string `env:"BOOTSTRAP_ADMIN_PASSWORD"`
	BootstrapAdminName string `env:"BOOTSTRAP_ADMIN_NAME"`
}

func (cfg *Config) validate() error {
	if _, err := net.ResolveTCPAddr("tcp", cfg.Addr); err != nil {
		return fmt.Errorf("invalid ADDR: %w", err)
	}

	if cfg.JWTKey == "" {
		return errors.New("JWT_KEY is required but not set")
	}

	if cfg.RateLimit <= 0 {
		return errors.New("RATE_LIMIT must be greater than 0")
	}

	if cfg.RateInterval <= 0 {
		return errors.New("RATE_INTERVAL must be greater than 0")
	}

	if cfg.DBPath == "" {
		return errors.New("DB_PATH is required")
	}

	if cfg.DBName == "" {
		return errors.New("DB_NAME is required")
	}

	if cfg.SchemaPath == "" {
		return errors.New("SCHEMA_PATH is required")
	}

	return nil
}

func MustLoadConfig() *Config {
	var cfg Config
	var envPath string

	flag.StringVar(&envPath, "config", "", "path to config file")
	flag.Parse()

	if envPath == "" {
		envPath = os.Getenv("CONFIG_PATH")
	}

	if envPath == "" {
		envPath = "config/config.env"
	}

	err := cleanenv.ReadConfig(envPath, &cfg)
	if err != nil {
		log.Fatalf("cannot read config file (%s): %v", envPath, err)
	}

	err = cfg.validate()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("config loaded from: %s", envPath)

	return &cfg
}
