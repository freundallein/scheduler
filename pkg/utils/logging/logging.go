package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

var logger = logrus.NewEntry(logrus.New())

// Fields describes auxiliary logged fields.
type Fields logrus.Fields

// Init used for a logger initialisation.
func Init(service, logLevel string) {
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
		"service": service,
	})
	logger.Debug("init_logger")
}

// WithFields used for adding auxiliary fields to event.
func WithFields(fields Fields) *logrus.Entry {
	return logger.WithFields(logrus.Fields(fields))
}

// Fatal logs fatal events.
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

// Error logs errors.
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Info logs events for information.
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Debug logs debug events.
func Debug(args ...interface{}) {
	logger.Debug(args...)
}
