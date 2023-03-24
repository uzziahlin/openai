package openai

import (
	"fmt"
	"io"
	"os"
)

const (
	LevelError = iota + 1
	LevelWarn
	LevelInfo
	LevelDebug
)

type Logger interface {
	Debugf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
}

type LeveledLogger struct {
	Level int

	stderrWriter io.Writer
	stdoutWriter io.Writer
}

func (l LeveledLogger) Debugf(format string, v ...interface{}) {
	if l.Level >= LevelDebug {
		fmt.Fprintf(l.stdout(), "[DEBUG] "+format+"\n", v...)
	}
}

func (l LeveledLogger) Errorf(format string, v ...interface{}) {
	if l.Level >= LevelError {
		fmt.Fprintf(l.stderr(), "[ERROR] "+format+"\n", v...)
	}
}

func (l LeveledLogger) Infof(format string, v ...interface{}) {
	if l.Level >= LevelInfo {
		fmt.Fprintf(l.stdout(), "[INFO] "+format+"\n", v...)
	}
}

func (l LeveledLogger) Warnf(format string, v ...interface{}) {
	if l.Level >= LevelWarn {
		fmt.Fprintf(l.stderr(), "[WARN] "+format+"\n", v...)
	}
}

func (l *LeveledLogger) stderr() io.Writer {
	if l.stderrWriter != nil {
		return l.stderrWriter
	}

	return os.Stderr
}

func (l *LeveledLogger) stdout() io.Writer {
	if l.stdoutWriter != nil {
		return l.stdoutWriter
	}

	return os.Stdout
}
