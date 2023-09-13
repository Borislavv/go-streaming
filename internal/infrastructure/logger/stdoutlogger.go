package logger

import (
	"os"
)

type StdOutLogger struct {
	*Logger
}

func NewStdOutLogger(buffer int) (logger *StdOutLogger, closeFunc func()) {
	l, closeFunc := NewLogger(os.Stdout, buffer)
	return &StdOutLogger{Logger: l}, closeFunc
}
