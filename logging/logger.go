package logging

import (
	"github.com/sirupsen/logrus"
)

// InitializeLogger initializes a logger with a default configuration
func InitializeLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	return logger
}