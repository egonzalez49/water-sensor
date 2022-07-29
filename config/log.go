package config

import "os"

type LogConfig struct {
	LogLevel string
}

func loadLogConfig() LogConfig {
	return LogConfig{
		LogLevel: os.Getenv("LOG_LEVEL"),
	}
}

// GetLogLevel returns the intended global log-level for the application.
// If the level is not specified, returns a level of 'info'.
func (l *LogConfig) GetLogLevel() []byte {
	if l.LogLevel != "" {
		return []byte(l.LogLevel)
	}

	return []byte("info")
}
