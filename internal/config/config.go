package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Env string

	HTTPHost        string
	HTTPPort        string
	ShutdownTimeout time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration

	//postressssss
	PGHost     string
	PGPort     string
	PGUser     string
	PGPassword string
	PGDatabase string
	PGSSLMode  string
	PGMaxConns int32
}

func Load() (Config, error) {
	cfg := Config{
		Env:             getEnv("APP_ENV", "dev"),
		HTTPHost:        getEnv("HTTP_HOST", "0.0.0.0"),
		HTTPPort:        getEnv("HTTP_PORT", "8080"),
		ShutdownTimeout: getDurationEnv("SHUTDOWN_TIMEOUT", 10*time.Second),
		ReadTimeout:     getDurationEnv("HTTP_READ_TIMEOUT", 10*time.Second),
		WriteTimeout:    getDurationEnv("HTTP_WRITE_TIMEOUT", 10*time.Second),
		IdleTimeout:     getDurationEnv("HTTP_IDLE_TIMEOUT", 60*time.Second),

		//post-gres xd
		PGHost:     getEnv("PG_HOST", "127.0.0.1"),
		PGPort:     getEnv("PG_PORT", "5432"),
		PGUser:     getEnv("PG_USER", "postgres"),
		PGPassword: getEnv("PG_PASSWORD", "postgres"),
		PGDatabase: getEnv("PG_DATABASE", "pos"),
		PGSSLMode:  getEnv("PG_SSLMODE", "disable"),
		PGMaxConns: getI32Env("PG_MAX_CONNS", 10),
	}

	cfg.Env = strings.ToLower(strings.TrimSpace(cfg.Env))
	if cfg.Env != "dev" && cfg.Env != "prod" {
		return Config{}, fmt.Errorf("variable de entorno inválida < APP_ENV: %q > debe ser <dev> o <prod>", cfg.Env)
	}

	if strings.TrimSpace(cfg.HTTPPort) == "" {
		return Config{}, fmt.Errorf("la variable HTTP_PORT es obligatoria")
	}

	if strings.TrimSpace(cfg.PGHost) == "" ||
		strings.TrimSpace(cfg.PGPort) == "" ||
		strings.TrimSpace(cfg.PGUser) == "" ||
		strings.TrimSpace(cfg.PGDatabase) == "" {
		return Config{}, fmt.Errorf("faltan variables de Postgres (PG_HOST/PG_PORT/PG_USER/PG_DATABASE)")
	}

	return cfg, nil
}

func getI32Env(s string, i int32) int32 {
	v, ok := os.LookupEnv(s)
	if !ok || strings.TrimSpace(v) == "" {
		return i
	}

	return parseI32OrDef(v, i)

}

func parseI32OrDef(v string, i int32) int32 {
	v = strings.TrimSpace(v)
	if v == "" {
		return i
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		return i
	}

	return int32(n)
}

func (c Config) GetHTTPAddr() string {
	return fmt.Sprintf("%s:%s", c.HTTPHost, c.HTTPPort)
}

func getDurationEnv(s string, duration time.Duration) time.Duration {
	v, ok := os.LookupEnv(s)
	if !ok || strings.TrimSpace(v) == "" {
		return duration
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		return duration
	}
	return d
}

func getEnv(key, val string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return val
}
