package log

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	ErrAlreadyClosed  = errors.New("Logger is already closed")
	ErrAlreadyBuilded = errors.New("Logger is already builded")
)

type Logger interface {
	Init(config string) error
	Write([]byte) (int, error)
	Sync() error
	Name() string
	GetLevel() Level
	GetStacktrackLevel() Level
	Close() error
}

type zapTeeLoggerConfig struct {
	zapcore.EncoderConfig
	Encoding string `json:"encoding" yaml:"encoding"`
}

type ZapTeeLogger struct {
	*zap.Logger
	lock    sync.Mutex
	loggers map[string]Logger
	closed  chan struct{}
	builded bool
}

func NewZapTeeLogger() *ZapTeeLogger {
	return &ZapTeeLogger{
		loggers: make(map[string]Logger),
		closed:  make(chan struct{}),
	}
}

func (l *ZapTeeLogger) AddLogger(name string, provider string, config string) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.builded {
		return ErrAlreadyBuilded
	}

	_, ok := l.loggers[name]
	if ok {
		delete(l.loggers, name)
	}

	logger, err := NewLogger(provider)
	if err != nil {
		return fmt.Errorf("NewLogger: %v", err)
	}
	if err := logger.Init(config); err != nil {
		return fmt.Errorf("Init: %v", err)
	}

	l.loggers[name] = logger
	return nil
}

func (l *ZapTeeLogger) Build(config string) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	var zapconfig zapTeeLoggerConfig
	zapconfig.EncoderConfig = zap.NewProductionEncoderConfig()
	zapconfig.Encoding = "console"
	zapconfig.EncodeTime = zapcore.RFC3339TimeEncoder

	if err := json.Unmarshal([]byte(config), &zapconfig); err != nil {
		return err
	}

	encoder := zapcore.NewConsoleEncoder
	if zapconfig.Encoding == "json" {
		encoder = zapcore.NewJSONEncoder
	}

	var stacktrackLevel Level = ErrorLevel
	var cores []zapcore.Core
	for _, logger := range l.loggers {
		syncer := zapcore.AddSync(logger)
		level := l.levelToZap(logger.GetLevel())
		enc := encoder(zapconfig.EncoderConfig)
		levelEnabler := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l <= level
		})
		core := zapcore.NewCore(enc, syncer, levelEnabler)
		cores = append(cores, core)

		if stacktrackLevel > logger.GetStacktrackLevel() {
			stacktrackLevel = logger.GetStacktrackLevel()
		}
	}

	zapStackTrackLevel := levelToZap(stacktrackLevel)
	levelEnbler := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= zapStackTrackLevel
	})

	l.Logger = zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddStacktrace(levelEnbler))
	l.builded = true
	return nil
}

func (l *ZapTeeLogger) levelToZap(level Level) zapcore.Level {
	return levelToZap(level)
}

func (l *ZapTeeLogger) Sync() error {
	if err := l.Logger.Sync(); err != nil {
		return err
	}

	l.lock.Lock()
	defer l.lock.Unlock()
	for name, logger := range l.loggers {
		if err := logger.Sync(); err != nil {
			return fmt.Errorf("Cannot sync logger(type: %s, name: %s): %v", logger.Name(), name, err)
		}
	}
	return nil
}

func (l *ZapTeeLogger) Close() error {
	select {
	case <-l.closed:
		return ErrAlreadyClosed
	default:
	}
	close(l.closed)

	l.lock.Lock()
	defer l.lock.Unlock()
	for name, logger := range l.loggers {
		if err := logger.Close(); err != nil {
			return fmt.Errorf("Cannot close logger(type: %s, name: %s): %v", logger.Name(), name, err)
		}
	}
	return nil
}
