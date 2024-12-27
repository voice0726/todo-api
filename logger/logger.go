package logger

import (
	zapPretty "github.com/thessem/zap-prettyconsole"
	"go.uber.org/zap"

	"github.com/voice0726/todo-app-api/config"
)

func NewLogger(c *config.Config) *zap.Logger {
	if c.IsProd {
		logger, _ := zap.NewProduction()
		return logger
	}
	return zapPretty.NewLogger(zap.DebugLevel)
}
