package log

import (
	"encoding/json"
	"fmt"

	"gopkg.in/natefinch/lumberjack.v2"
)

type FileLogger struct {
	lumberjack.Logger
	Level           string `json:"level"`
	StacktrackLevel string `json:"stacktracklevel"`
	level           Level
	stacktracklevel Level
}

func (l *FileLogger) Init(config string) error {
	if err := json.Unmarshal([]byte(config), l); err != nil {
		return err
	}

	if len(l.Filename) == 0 {
		return fmt.Errorf("config must have filename")
	}

	l.level = NameToLevel(l.Level)
	l.stacktracklevel = NameToLevel(l.StacktrackLevel)
	return nil
}

func (l *FileLogger) Write(p []byte) (int, error) {
	return l.Logger.Write(p)
}

func (l *FileLogger) Sync() error {
	return nil
}

func (l *FileLogger) Name() string {
	return "file"
}

func (l *FileLogger) GetLevel() Level {
	return l.level
}

func (l *FileLogger) GetStacktrackLevel() Level {
	return l.stacktracklevel
}

func (l *FileLogger) Close() error {
	return l.Logger.Close()
}

func NewFileLogger() Logger {
	return &FileLogger{}
}

func init() {
	Register("file", NewFileLogger)
}
