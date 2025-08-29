package logger

import (
	"context"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// Logger is a custom logger wrapper
type Logger struct {
	*logrus.Logger
}

// Config holds logger configuration
type Config struct {
	Level       string `json:"level"`
	Format      string `json:"format"` // "json" or "text"
	Output      string `json:"output"` // "stdout", "stderr", or file path
	ServiceName string `json:"service_name"`
}

// NewLogger creates a new configured logger instance
func NewLogger(config Config) *Logger {
	logger := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// Set output format
	switch strings.ToLower(config.Format) {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "@timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
				logrus.FieldKeyFunc:  "function",
				logrus.FieldKeyFile:  "file",
			},
		})
	default:
		logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		})
	}

	// Set output destination
	switch strings.ToLower(config.Output) {
	case "stderr":
		logger.SetOutput(os.Stderr)
	case "stdout", "":
		logger.SetOutput(os.Stdout)
	default:
		// Assume it's a file path
		if file, err := os.OpenFile(config.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err == nil {
			logger.SetOutput(file)
		} else {
			logger.SetOutput(os.Stdout)
			logger.WithError(err).Warn("Failed to open log file, using stdout")
		}
	}

	// Add service name as default field
	if config.ServiceName != "" {
		logger = logger.WithField("service", config.ServiceName).Logger
	}

	return &Logger{Logger: logger}
}

// WithContext adds context fields to the logger
func (l *Logger) WithContext(ctx context.Context) *logrus.Entry {
	entry := l.WithFields(logrus.Fields{})

	// Extract request ID from context if available
	if requestID := ctx.Value("request_id"); requestID != nil {
		entry = entry.WithField("request_id", requestID)
	}

	// Extract user ID from context if available
	if userID := ctx.Value("user_id"); userID != nil {
		entry = entry.WithField("user_id", userID)
	}

	// Extract merchant ID from context if available
	if merchantID := ctx.Value("merchant_id"); merchantID != nil {
		entry = entry.WithField("merchant_id", merchantID)
	}

	return entry
}

// WithRequestID adds request ID to the logger
func (l *Logger) WithRequestID(requestID string) *logrus.Entry {
	return l.WithField("request_id", requestID)
}

// WithUserID adds user ID to the logger
func (l *Logger) WithUserID(userID string) *logrus.Entry {
	return l.WithField("user_id", userID)
}

// WithMerchantID adds merchant ID to the logger
func (l *Logger) WithMerchantID(merchantID string) *logrus.Entry {
	return l.WithField("merchant_id", merchantID)
}

// WithError adds error information to the logger
func (l *Logger) WithError(err error) *logrus.Entry {
	return l.Logger.WithError(err)
}

// WithFields adds multiple fields to the logger
func (l *Logger) WithFields(fields map[string]interface{}) *logrus.Entry {
	return l.Logger.WithFields(logrus.Fields(fields))
}

// Database logs database-related events
func (l *Logger) Database() *logrus.Entry {
	return l.WithField("component", "database")
}

// HTTP logs HTTP-related events
func (l *Logger) HTTP() *logrus.Entry {
	return l.WithField("component", "http")
}

// Auth logs authentication-related events
func (l *Logger) Auth() *logrus.Entry {
	return l.WithField("component", "auth")
}

// Business logs business logic events
func (l *Logger) Business() *logrus.Entry {
	return l.WithField("component", "business")
}

// External logs external service interactions
func (l *Logger) External() *logrus.Entry {
	return l.WithField("component", "external")
}

// DefaultConfig returns a default logger configuration
func DefaultConfig(serviceName string) Config {
	return Config{
		Level:       "info",
		Format:      "json",
		Output:      "stdout",
		ServiceName: serviceName,
	}
}

// Global logger instance (can be used for simple logging)
var globalLogger *Logger

// Initialize sets up the global logger
func Initialize(config Config) {
	globalLogger = NewLogger(config)
}

// GetGlobalLogger returns the global logger instance
func GetGlobalLogger() *Logger {
	if globalLogger == nil {
		globalLogger = NewLogger(DefaultConfig("unified-commerce"))
	}
	return globalLogger
}

// Convenience functions using the global logger
func Debug(args ...interface{}) {
	GetGlobalLogger().Debug(args...)
}

func Info(args ...interface{}) {
	GetGlobalLogger().Info(args...)
}

func Warn(args ...interface{}) {
	GetGlobalLogger().Warn(args...)
}

func Error(args ...interface{}) {
	GetGlobalLogger().Error(args...)
}

func Fatal(args ...interface{}) {
	GetGlobalLogger().Fatal(args...)
}

func WithFields(fields map[string]interface{}) *logrus.Entry {
	return GetGlobalLogger().WithFields(fields)
}

func WithError(err error) *logrus.Entry {
	return GetGlobalLogger().WithError(err)
}
