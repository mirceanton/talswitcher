package main

import (
	"os"
	"strings"

	"github.com/mirceanton/talswitcher/cmd"
	log "github.com/sirupsen/logrus"
)

func SetLogLevel() {
	logLevel := os.Getenv("TALSWITCHER_LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	level, err := log.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		log.Fatalf("Invalid log level: %v", err)
	}
	log.SetLevel(level)
}

func SetLogFormat() {
	logFormat := os.Getenv("TALSWITCHER_LOG_FORMAT")
	if logFormat == "" {
		logFormat = "text"
	}

	switch strings.ToLower(logFormat) {
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	case "text":
		log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	default:
		log.Fatalf("Invalid logformat: %v", logFormat)
	}
}

func init() {
	SetLogLevel()
	SetLogFormat()
}

func main() {
	cmd.Execute()
}
