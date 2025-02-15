package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/lmittmann/tint"
)

type Config struct {
	LogLevel     string        `default:"info" split_words:"true"`
	DatabasePath string        `default:"interruptions.db" split_words:"true"`
	BindAddr     string        `default:":8080" split_words:"true"`
	PollInterval time.Duration `default:"5m" split_words:"true"`
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		slog.Warn("Failed to load .env file", "err", err)
	}

	var cfg Config
	err = envconfig.Process("interruption_spotter", &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	level, validLevel := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}[strings.ToLower(cfg.LogLevel)]
	if !validLevel {
		return nil, fmt.Errorf("invalid log level: %s", cfg.LogLevel)
	}

	slog.SetDefault(slog.New(
		tint.NewHandler(
			os.Stderr,
			&tint.Options{
				TimeFormat: time.Kitchen,
				Level:      level,
			},
		),
	))

	return &cfg, nil
}
