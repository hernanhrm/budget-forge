package logger

import (
	"context"

	shared_domain "github.com/hernanhrm/budget-forge/pkg/shared_domain"
)

type NoopLogger struct{}

func NewNoop() shared_domain.Logger {
	return NoopLogger{}
}

func (n NoopLogger) Debug(_ string, _ ...interface{}) {}

func (n NoopLogger) Info(_ string, _ ...interface{}) {}

func (n NoopLogger) Warn(_ string, _ ...interface{}) {}

func (n NoopLogger) Error(_ string, _ ...interface{}) {}

func (n NoopLogger) With(_ ...interface{}) shared_domain.Logger {
	return n
}

func (n NoopLogger) WithContext(_ context.Context) shared_domain.Logger {
	return n
}
