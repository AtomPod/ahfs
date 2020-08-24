package log

import (
	"encoding/json"
	"os"
)

type ConsoleLogger struct {
	Level           string `json:"level" yaml:"level"`
	StacktrackLevel string `json:"stackTrackLevel" yaml:"stackTrackLevel"`
	Stderr          bool   `json:"stderr" yaml:"stderr"`
	level           Level
	stacktracklevel Level
	out             *os.File
}

func NewConsoleLogger() Logger {
	return &ConsoleLogger{}
}

func (l *ConsoleLogger) Init(config string) error {
	if err := json.Unmarshal([]byte(config), l); err != nil {
		return err
	}

	l.out = os.Stdout
	if l.Stderr {
		l.out = os.Stderr
	}

	l.level = NameToLevel(l.Level)
	if len(l.StacktrackLevel) == 0 {
		l.stacktracklevel = ErrorLevel
	} else {
		l.stacktracklevel = NameToLevel(l.StacktrackLevel)
	}

	return nil
}

func (l *ConsoleLogger) Write(p []byte) (int, error) {
	return l.out.Write(p)
}

func (l *ConsoleLogger) Sync() error {
	return l.out.Sync()
}

func (l *ConsoleLogger) Name() string {
	return "console"
}

func (l *ConsoleLogger) GetLevel() Level {
	return l.level
}

func (l *ConsoleLogger) GetStacktrackLevel() Level {
	return l.stacktracklevel
}

func (l *ConsoleLogger) Close() error {
	return nil
}

func init() {
	Register("console", NewConsoleLogger)
}
