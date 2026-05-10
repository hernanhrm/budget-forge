package logger

import (
	"context"
	"io"
	"log/slog"
	"os"

	shared_domain "github.com/hernanhrm/budget-forge/pkg/shared_domain"
)

type SlogAdapter struct {
	logger *slog.Logger
}

func NewSlogAdapter(config Config) shared_domain.Logger {
	return newSlogAdapter(config, os.Stdout)
}

func newSlogAdapter(config Config, writer io.Writer) shared_domain.Logger {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: convertLevel(config.Level),
	}

	switch config.Format {
	case FormatJSON:
		handler = slog.NewJSONHandler(writer, opts)
	case FormatText:
		handler = slog.NewTextHandler(writer, opts)
	default:
		handler = slog.NewTextHandler(writer, opts)
	}

	return SlogAdapter{
		logger: slog.New(handler),
	}
}

func (s SlogAdapter) Debug(msg string, keysAndValues ...any) {
	s.logger.Debug(msg, keysAndValues...)
}

func (s SlogAdapter) Info(msg string, keysAndValues ...any) {
	s.logger.Info(msg, keysAndValues...)
}

func (s SlogAdapter) Warn(msg string, keysAndValues ...any) {
	s.logger.Warn(msg, keysAndValues...)
}

func (s SlogAdapter) Error(msg string, keysAndValues ...any) {
	s.logger.Error(msg, keysAndValues...)
}

func (s SlogAdapter) With(keysAndValues ...any) shared_domain.Logger {
	return SlogAdapter{
		logger: s.logger.With(keysAndValues...),
	}
}

func (s SlogAdapter) WithContext(_ context.Context) shared_domain.Logger {
	return SlogAdapter{
		logger: s.logger.With(),
	}
}

func convertLevel(level Level) slog.Level {
	switch level {
	case LevelDebug:
		return slog.LevelDebug
	case LevelInfo:
		return slog.LevelInfo
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func NewDevelopment() shared_domain.Logger {
	return NewSlogAdapter(Config{
		Level:  LevelDebug,
		Format: FormatText,
	})
}

func NewProduction() shared_domain.Logger {
	return NewSlogAdapter(Config{
		Level:  LevelInfo,
		Format: FormatJSON,
	})
}
