package log

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

func levelToZap(level Level) zapcore.Level {
	switch level {
	case DebugLevel:
		return zap.DebugLevel
	case InfoLevel:
		return zap.InfoLevel
	case WarnLevel:
		return zap.WarnLevel
	case ErrorLevel:
		return zap.ErrorLevel
	case FatalLevel:
		return zap.FatalLevel
	}
	return zap.InfoLevel
}

func NameToLevel(name string) Level {
	name = strings.ToLower(name)
	switch name {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	}
	return InfoLevel
}
