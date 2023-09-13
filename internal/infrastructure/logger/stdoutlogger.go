package logger

import (
	"os"
)

type StdOutLogger struct {
	*Logger
}

func NewStdOutLogger(errBuff int, reqBuff int) (logger *StdOutLogger, closeFunc func()) {
	l, closeFunc := NewLogger(os.Stdout, errBuff, reqBuff)
	return &StdOutLogger{Logger: l}, closeFunc
}
