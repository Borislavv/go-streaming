package logger

import (
	"context"
	"os"
)

type StdErr struct {
	*abstract
}

func NewStdErr(ctx context.Context, errBuff int, reqBuff int) (logger *StdErr, closeFunc func()) {
	abstractLogger, closeFunc := newAbstractLogger(ctx, os.Stderr, errBuff, reqBuff)
	return &StdErr{abstract: abstractLogger}, closeFunc
}
