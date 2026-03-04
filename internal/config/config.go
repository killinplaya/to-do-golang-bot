package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken string
	LogLevel slog.Level
}

func Load() (Config, error) {
	if err := loadDotEnv(); err != nil {
		return Config{}, err
	}

	cfg := Config{
		BotToken: strings.TrimSpace(os.Getenv("BOT_TOKEN")),
		LogLevel: slog.LevelInfo,
	}

	if cfg.BotToken == "" {
		return Config{}, errors.New("BOT_TOKEN is required")
	}

	logLevelRaw := strings.TrimSpace(os.Getenv("LOG_LEVEL"))
	if logLevelRaw == "" {
		return cfg, nil
	}

	level, err := parseLogLevel(logLevelRaw)
	if err != nil {
		return Config{}, err
	}
	cfg.LogLevel = level

	return cfg, nil
}

func parseLogLevel(v string) (slog.Level, error) {
	switch strings.ToLower(v) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("invalid LOG_LEVEL: %q", v)
	}
}

func loadDotEnv() error {
	if err := godotenv.Load(); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("load .env: %w", err)
	}

	return nil
}
