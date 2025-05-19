package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v10"
)

// Config общий конфиг
type Config struct {
	Env        string `env:"ENV" envDefault:"prod"`
	Postgres   Postgres
	HTTPServer HTTPServer
}

type HTTPServer struct {
	Host        string `env:"HTTP_HOST" envDefault:"localhost"`
	Port        string `env:"HTTP_PORT" envDefault:"8080"`
	Address     string
	Timeout     time.Duration `env:"HTTP_TIMEOUT" envDefault:"5s"`
	IdleTimeout time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"60s"`
	User        string        `env:"HTTP_USER" envDefault:"user"`
	Password    string        `env:"HTTP_PASSWORD" envDefault:"password"`
}

type Postgres struct {
	Host     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	Port     string `env:"POSTGRES_PORT" envDefault:"5432"`
	User     string `env:"POSTGRES_USER" envDefault:"user"`
	Password string `env:"POSTGRES_PASSWORD" envDefault:"password"`
	Db       string `env:"POSTGRES_DB" envDefault:"postgres"`
	SslMode  string `env:"POSTGRES_SSL_MODE" envDefault:"disable"`
	DSN      string
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("loading config from env is failed: %w", err)
	}
	buildAddress(&cfg.HTTPServer)
	buildDSN(&cfg.Postgres)

	return cfg, nil
}

func buildAddress(httpserver *HTTPServer) {
	httpserver.Address = fmt.Sprintf("%s:%s", httpserver.Host, httpserver.Port)
}

func buildDSN(p *Postgres) {
	p.DSN = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		p.User, p.Password, p.Host, p.Port, p.Db, p.SslMode)
}
