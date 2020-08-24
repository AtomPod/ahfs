package log

import (
	"errors"
	"fmt"
	"sync"

	"go.uber.org/zap"
)

type Provider func() Logger

var (
	providersMux  sync.RWMutex
	providers     = map[string]Provider{}
	defaultLogger *ZapTeeLogger
)

func Register(name string, p Provider) {
	if p == nil {
		panic("Logger: register provider is nil")
	}

	providersMux.Lock()
	defer providersMux.Unlock()

	if _, ok := providers[name]; ok {
		panic("Logger: register provider twice for provider: " + name)
	}
	providers[name] = p
}

func NewLogger(providerName string) (Logger, error) {
	providersMux.RLock()
	defer providersMux.RUnlock()

	provider, ok := providers[providerName]
	if !ok {
		return nil, fmt.Errorf("Logger: cannot found provider (%s)", providerName)
	}

	return provider(), nil
}

func ZapLogger() *zap.Logger {
	return defaultLogger.Logger
}

func Init() {
	if defaultLogger == nil {
		defaultLogger = NewZapTeeLogger()
	}
}

func Sync() {
	if defaultLogger != nil {
		if err := defaultLogger.Sync(); err != nil {
			Error("Cannot sync logger", zap.Error(err))
		}
	}
}

func AddLogger(name string, provider string, config string) error {
	Init()
	return defaultLogger.AddLogger(name, provider, config)
}

func New(config string) error {
	if defaultLogger == nil {
		return errors.New("Logger: logger is nil, please run Init()")
	}
	return defaultLogger.Build(config)
}

func Reset() error {
	if defaultLogger != nil {
		if err := defaultLogger.Close(); err != nil {
			return err
		}
	}
	defaultLogger = NewZapTeeLogger()
	return nil
}

func SetDefault(l *ZapTeeLogger) {
	if defaultLogger != nil {
		if err := defaultLogger.Close(); err != nil {
			l.Warn("Cannot close logger", zap.Error(err))
		}
	}
	defaultLogger = l
}

func Debug(msg string, fields ...zap.Field) {
	if defaultLogger == nil {
		return
	}
	defaultLogger.WithOptions(zap.AddCallerSkip(1)).Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	if defaultLogger == nil {
		return
	}
	defaultLogger.WithOptions(zap.AddCallerSkip(1)).Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	if defaultLogger == nil {
		return
	}
	defaultLogger.WithOptions(zap.AddCallerSkip(1)).Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	if defaultLogger == nil {
		return
	}
	defaultLogger.WithOptions(zap.AddCallerSkip(1)).Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	if defaultLogger == nil {
		return
	}
	defaultLogger.WithOptions(zap.AddCallerSkip(1)).Fatal(msg, fields...)
}

func DebugWithSkip(skip int, msg string, fields ...zap.Field) {
	if defaultLogger == nil {
		return
	}
	defaultLogger.WithOptions(zap.AddCallerSkip(1+skip)).Debug(msg, fields...)
}

func InfoWithSkip(skip int, msg string, fields ...zap.Field) {
	if defaultLogger == nil {
		return
	}
	defaultLogger.WithOptions(zap.AddCallerSkip(1+skip)).Info(msg, fields...)
}

func WarnWithSkip(skip int, msg string, fields ...zap.Field) {
	if defaultLogger == nil {
		return
	}
	defaultLogger.WithOptions(zap.AddCallerSkip(1+skip)).Warn(msg, fields...)
}

func ErrorWithSkip(skip int, msg string, fields ...zap.Field) {
	if defaultLogger == nil {
		return
	}
	defaultLogger.WithOptions(zap.AddCallerSkip(1+skip)).Error(msg, fields...)
}

func FatalWithSkip(skip int, msg string, fields ...zap.Field) {
	if defaultLogger == nil {
		return
	}
	defaultLogger.WithOptions(zap.AddCallerSkip(1+skip)).Fatal(msg, fields...)
}
