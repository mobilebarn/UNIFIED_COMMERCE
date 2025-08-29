package database

import (
	"gorm.io/gorm/logger"
)

// NewGormLogger creates a new GORM logger instance
func NewGormLogger() logger.Interface {
	return logger.Default.LogMode(logger.Info)
}
