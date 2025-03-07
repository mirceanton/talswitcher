package log

import (
	"os"

	"github.com/mirceanton/talswitcher/pkg/types"
	log "github.com/sirupsen/logrus"
)

// Setup initializes the logger with the provided configuration
func Setup(config types.Config) {
	setLogLevel(config.LogLevel)
	setLogFormat(config.LogFormat)
	log.Debugf("Logger initialized with level: %s, format: %s", config.LogLevel, config.LogFormat)
}

func setLogLevel(logLevel string) {
	if logLevel == "" {
		logLevel = os.Getenv("TALSWITCHER_LOG_LEVEL")
		if logLevel == "" {
			logLevel = "info"
		}
	}
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Fatalf("Invalid log level: %v", err)
	}
	log.SetLevel(level)
}

func setLogFormat(logFormat string) {
	if logFormat == "" {
		logFormat = os.Getenv("TALSWITCHER_LOG_FORMAT")
		if logFormat == "" {
			logFormat = "text"
		}
	}
	switch logFormat {
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	case "text":
		log.SetFormatter(&log.TextFormatter{})
	default:
		log.Fatalf("Invalid log format: %s", logFormat)
	}
}
