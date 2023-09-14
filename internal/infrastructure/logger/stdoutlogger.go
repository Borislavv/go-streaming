package logger

import (
	"context"
	"os"
)

type StdOutLogger struct {
	*Logger
}

func NewStdOutLogger(ctx context.Context, errBuff int, reqBuff int) (logger *StdOutLogger, closeFunc func()) {
	l, closeFunc := NewLogger(ctx, os.Stdout, errBuff, reqBuff)
	return &StdOutLogger{Logger: l}, closeFunc
}
