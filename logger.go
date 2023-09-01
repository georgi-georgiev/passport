package passport

import (
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() *zap.Logger {
	logger, _ := zap.NewDevelopment()

	sentryOptions := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.RegisterHooks(core, func(entry zapcore.Entry) error {
			if entry.Level == zapcore.ErrorLevel {
				sentry.CaptureMessage(entry.Message)
			}
			return nil
		})
	})

	logger.WithOptions(sentryOptions)

	return logger
}
