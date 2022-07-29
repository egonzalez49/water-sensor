package logging

import (
	"github.com/egonzalez49/water-sensor/config"
	"go.uber.org/zap"
)

type Logger struct {
	*zap.SugaredLogger
	config zap.Config
}

func NewLogger(cfg *config.Config) (*Logger, error) {
	logCfg := zap.NewProductionConfig()

	logLevel := cfg.Log.GetLogLevel()
	logCfg.Level.UnmarshalText(logLevel)

	logger, err := logCfg.Build()
	if err != nil {
		return nil, err
	}

	sugaredLogger := logger.Sugar()

	return &Logger{sugaredLogger, logCfg}, nil
}

func (l *Logger) Shutdown() error {
	return l.Sync()
}
