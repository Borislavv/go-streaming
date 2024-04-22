package logger

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/infrastructure/helper"
	"os"
	"time"
)

type File struct {
	*abstract
}

func NewFile(ctx context.Context, errBuff int, reqBuff int) (logger *File, closeFunc func(), err error) {
	logsDir, err := helper.LogsDir()
	if err != nil {
		return nil, nil, err
	}

	_, err = os.Open(logsDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, nil, err
		}
		if err = os.MkdirAll(logsDir, os.ModeDir); err != nil {
			return nil, nil, err
		}
	}

	logfile, err := os.Create(logsDir + time.Now().Format("2006_01_02") + ".log")
	if err != nil {
		return nil, nil, err
	}

	abstractLogger, closeFunc := newAbstractLogger(ctx, logfile, errBuff, reqBuff)

	return &File{abstract: abstractLogger}, closeFunc, nil
}
