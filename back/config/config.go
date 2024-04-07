package config

import (
	"flag"
	"fmt"
	"log/slog"
	"strings"

	"github.com/caarlos0/env/v10"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	ServerPort string `env:"SERVER_PORT" envDefault:"8080"`
	ServerURL  string `env:"SERVER_URL" envDefault:"localhost"`

	logLevelString string     `env:"LOG_LEVEL" envDefault:"info"`
	LogLevel       slog.Level // Parse from logLevelString

	PSQLDSN      string `env:"PSQL_DSN" envDefault:""` // If PSQL_DSN is set, other PSQL_* variables will be ignored
	PSQLHost     string `env:"PSQL_HOST" envDefault:"localhost"`
	PSQLPort     string `env:"PSQL_PORT" envDefault:"5432"`
	PSQLDatabase string `env:"PSQL_DATABASE" envDefault:"ac_hacking"`
	PSQLUser     string `env:"PSQL_USER" envDefault:"postgres"`
	PSQLPassword string `env:"PSQL_PASSWORD" envDefault:"postgres"`
	PSQLSSLMode  string `env:"PSQL_SSL_MODE" envDefault:"disable"`
	PSQLTimeZone string `env:"PSQL_TIMEZONE" envDefault:"Asia/Tokyo"`

	FrontURL string `env:"FRONT_URL" envDefault:"http://localhost:3000"`

	ProxyPassword string `env:"PROXY_PASSWORD" envDefault:""`
}

func Load() (Config, error) {
	serverport := flag.String("port", "8080", "server port")
	flag.Parse()

	cfg := Config{}
	if serverport != nil {
		cfg.ServerPort = *serverport
	}

	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}

	cfg.LogLevel = parseLogLevel(cfg.logLevelString)
	cfg = makePSQLDSN(cfg)

	return cfg, nil
}

func parseLogLevel(str string) slog.Level {
	switch strings.ToLower(str) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		slog.Warn("Invalid log level, use default level: info")
		return slog.LevelInfo
	}
}

func makePSQLDSN(config Config) Config {
	if config.PSQLDSN != "" {
		return config
	}

	config.PSQLDSN = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s", config.PSQLHost, config.PSQLPort, config.PSQLUser, config.PSQLPassword, config.PSQLDatabase, config.PSQLSSLMode, config.PSQLTimeZone)

	return config
}
