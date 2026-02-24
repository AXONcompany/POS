package config

import (
	"fmt"
	"net/url"
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

	DBMaxOpenConns int32
	DBMaxIdleConns int32
	DBConnMaxLife  time.Duration
	JWTSecret      string
}

func Load() (Config, error) {
	cfg := Config{
		Env:             getEnv("APP_ENV", "dev"),
		HTTPHost:        getEnv("HTTP_HOST", "0.0.0.0"),
		HTTPPort:        getEnv("HTTP_PORT", "8081"),
		ShutdownTimeout: getDurationEnv("SHUTDOWN_TIMEOUT", 10*time.Second),
		ReadTimeout:     getDurationEnv("HTTP_READ_TIMEOUT", 10*time.Second),
		WriteTimeout:    getDurationEnv("HTTP_WRITE_TIMEOUT", 10*time.Second),
		IdleTimeout:     getDurationEnv("HTTP_IDLE_TIMEOUT", 60*time.Second),

		PGHost:     getEnvFirst("PG_HOST", "DB_HOST", "POSTGRES_HOST", "127.0.0.1"),
		PGPort:     getEnvFirst("PG_PORT", "DB_PORT", "POSTGRES_PORT", "5433"),
		PGUser:     getEnvFirst("PG_USER", "DB_USER", "POSTGRES_USER", "postgres"),
		PGPassword: getEnvFirst("PG_PASSWORD", "DB_PASSWORD", "POSTGRES_PASSWORD", "postgres"),
		PGDatabase: getEnvFirst("PG_DATABASE", "DB_NAME", "POSTGRES_DB", "pos"),
		PGSSLMode:  getEnvFirst("PG_SSLMODE", "DB_SSLMODE", "POSTGRES_SSLMODE", "disable"),
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

func (c Config) GetPostgresDSN() string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.PGUser, c.PGPassword),
		Host:   fmt.Sprintf("%s:%s", c.PGHost, c.PGPort),
		Path:   c.PGDatabase,
	}

	q := u.Query()
	q.Set("sslmode", c.PGSSLMode)
	u.RawQuery = q.Encode()
	return u.String()
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

func getEnvFirst(keys ...string) string {
	if len(keys) == 0 {
		return ""
	}

	def := keys[len(keys)-1]
	for _, k := range keys[:len(keys)-1] {
		if v, ok := os.LookupEnv(k); ok {
			v = strings.TrimSpace(v)
			if v != "" {
				return v
			}
		}
	}
	return def
}
