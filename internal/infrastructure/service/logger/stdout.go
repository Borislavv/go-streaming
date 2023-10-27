package logger

import (
	"context"
	"os"
)

type StdOut struct {
	*abstract
}

func NewStdOut(ctx context.Context, errBuff int, reqBuff int) (logger *StdOut, closeFunc func()) {
	abstractLogger, closeFunc := newAbstractLogger(ctx, os.Stdout, errBuff, reqBuff)
	return &StdOut{abstract: abstractLogger}, closeFunc
}
