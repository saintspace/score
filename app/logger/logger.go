package logger

import (
	"fmt"

	"go.uber.org/zap"
)

type Logger struct {
	logger        *zap.Logger
	sugaredLogger *zap.SugaredLogger
}

func New(debugEnabled bool) (*Logger, error) {
	var logger *zap.Logger
	var err error
	if debugEnabled {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		return nil, fmt.Errorf("error initializing logger => %v", err.Error())
	}
	return &Logger{
		logger:        logger,
		sugaredLogger: logger.Sugar(),
	}, nil
}

func (s *Logger) InfoWithContext(message string, keysAndValues ...interface{}) {
	s.sugaredLogger.Infow(message, keysAndValues...)
}

func (s *Logger) ErrorWithContext(message string, keysAndValues ...interface{}) {
	s.sugaredLogger.Errorw(message, keysAndValues...)
}

func (s *Logger) DebugWithContext(message string, keysAndValues ...interface{}) {
	s.sugaredLogger.Debugw(message, keysAndValues...)
}
