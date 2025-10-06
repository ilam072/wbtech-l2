package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"net"
	"time"
)

type Config struct {
	DBConfig     DBConfig
	ServerConfig ServerConfig
	JWTConfig    JWTConfig
}

type DBConfig struct {
	PgUser     string `env:"PGUSER"`
	PgPassword string `env:"PGPASSWORD"`
	PgHost     string `env:"PGHOST"`
	PgPort     uint16 `env:"PGPORT"`
	PgDatabase string `env:"PGDATABASE"`
	PgSSLMode  string `env:"PGSSLMODE"`
}

type ServerConfig struct {
	HTTPPort string `env:"HTTP_PORT"`
}

type JWTConfig struct {
	SecretKey string `env:"SECRET"`
	TokenTTL  time.Duration
}

func (s *ServerConfig) Address() string {
	return net.JoinHostPort("localhost", s.HTTPPort)

}

func New() *Config {
	cfg := &Config{}

	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	if err := env.Parse(cfg); err != nil {
		panic(err)
	}

	return cfg
}
