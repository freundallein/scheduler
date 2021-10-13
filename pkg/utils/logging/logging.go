package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

var logger = logrus.NewEntry(logrus.New())

// Fields ...
type Fields logrus.Fields

// Init ...
func Init(module, logLevel string) {
	customFormatter := &logrus.TextFormatter{}
	customFormatter.TimestampFormat = timeFormat
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)
	logrus.SetOutput(os.Stdout)
	switch logLevel {
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
	logger = logrus.WithFields(logrus.Fields{
		"module": module,
	})
	logger.Debug("init_logger")
}

// WithFields ...
func WithFields(fields Fields) *logrus.Entry {
	return logger.WithFields(logrus.Fields(fields))
}

// Fatal ...
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

// Error ...
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Info ...
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Debug ...
func Debug(args ...interface{}) {
	logger.Debug(args...)
}
