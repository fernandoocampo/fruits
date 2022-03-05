package loggers

import (
	"io"
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

// Level log level type
type Level int

// Fields define fields type to add in logs.
type Fields map[string]interface{}

const (
	// Debug log level
	Debug Level = 1 + iota
	// Info log level
	Info
	// Warn log level
	Warn
	// Error log level
	Error
)

// artifactField is the name of the artifact field in the logs.
const artifactField = "artifact"

// Logger contains a logrus logger.
type Logger struct {
	// level log level
	level Level
	// artifact is the name of the artifact that print the log.
	artifact string
	// logger is the logrus logger.
	logger *logrus.Logger
}

// NewLoggerWithStdout create a logrus logger with stdout output.
func NewLoggerWithStdout(artifactName string, level Level) *Logger {
	return NewLogger(artifactName, level, os.Stdout)
}

// NewLogger create a logrus logger.
func NewLogger(artifactName string, level Level, output io.Writer) *Logger {
	logger := logrus.New()
	logger.SetOutput(output)
	logger.SetLevel(getLoggerLevel(level))
	logger.SetFormatter(&logrus.JSONFormatter{})
	newLogrus := Logger{
		level:    level,
		logger:   logger,
		artifact: artifactName,
	}
	return &newLogrus
}

// NewBasicLogger create basic logger
func NewBasicLogger(output io.Writer) *log.Logger {
	newLogger := log.New(output, "", 0)
	return newLogger
}

func (l *Logger) Debug(msg string, fields Fields) {
	if l.level&Debug != 0 {
		l.logger.WithField(artifactField, l.artifact).
			WithFields(logrus.Fields(fields)).
			Debug(msg)
	}
}

func (l *Logger) Info(msg string, fields Fields) {
	if l.level&(Debug|Info) != 0 {
		l.logger.WithField(artifactField, l.artifact).
			WithFields(logrus.Fields(fields)).
			Info(msg)
	}
}

func (l *Logger) Warn(msg string, fields Fields) {
	if l.level&(Debug|Info|Warn) != 0 {
		l.logger.WithField(artifactField, l.artifact).
			WithFields(logrus.Fields(fields)).
			Warn(msg)
	}
}

func (l *Logger) Error(msg string, fields Fields) {
	if l.level&(Debug|Info|Warn|Error) != 0 {
		l.logger.WithField(artifactField, l.artifact).
			WithFields(logrus.Fields(fields)).
			Error(msg)
	}
}

// SetLoggerLevel set the level of log
func (l *Logger) SetLoggerLevel(newLevel Level) {
	if newLevel < 0 {
		l.level = Debug
		return
	}
	if newLevel > Error {
		l.level = Error
		return
	}
	l.level = newLevel
}

// getLoggerLevel get the logger level based on allowed levels.
func getLoggerLevel(level Level) logrus.Level {
	switch level {
	case Debug:
		return logrus.DebugLevel
	case Info:
		return logrus.InfoLevel
	case Warn:
		return logrus.WarnLevel
	case Error:
		return logrus.ErrorLevel
	default:
		return logrus.DebugLevel
	}
}
