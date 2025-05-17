package config

import (
	"fmt"
	"github.com/caarlos0/env/v10"
	"time"
)

// Config общий конфиг
type Config struct {
	Environment string `env:"ENV" envDefault:"local"`
	Postgres    Postgres
	HTTPServer  HTTPServer
}

type HTTPServer struct {
	Address     string        `yaml:"address" envDefault:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" envDefault:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" envDefault:"60s"`
	User        string        `yaml:"user" envRequired:"true"`
	Password    string        `yaml:"password" envRequired:"true" env:"HTTP_SERVER_PASSWORD"`
}

type Postgres struct {
	Host     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	Port     string `env:"POSTGRES_PORT" envDefault:"5432"`
	User     string `env:"POSTGRES_USER" envDefault:"root"`
	Password string `env:"POSTGRES_PASSWORD" envDefault:"password"`
	Db       string `env:"POSTGRES_DB" envDefault:"postgres"`
	SslMode  string `env:"POSTGRES_SSL_MODE" envDefault:"disable"`
	DSN      string `env:"POSTGRES_DSN"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("loading config from env is failed: %w", err)
	}
	buildDSN(&cfg.Postgres)

	return cfg, nil
}

func buildDSN(p *Postgres) {
	p.DSN = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		p.User, p.Password, p.Host, p.Port, p.Db, p.SslMode)
}
