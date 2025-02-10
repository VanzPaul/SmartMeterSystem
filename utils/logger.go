package utils

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

// NewLogger creates a new logger instance based on the debug flag.
func NewLogger(debug bool) (*zap.Logger, error) {
	var Logger *zap.Logger
	var err error

	if debug {
		Logger, err = zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
		Logger.Info("Logger initialized in development mode.")
	} else {
		config := zap.NewProductionConfig()
		config.Level.SetLevel(zap.InfoLevel)
		Logger, err = config.Build()
		if err != nil {
			return nil, err
		}
		Logger.Info("Logger initialized in production mode.")
	}

	return Logger, nil
}
