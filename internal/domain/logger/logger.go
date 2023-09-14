package logger

import (
	"context"
	"io"
)

type Logger interface {
	Log(err error)
	LogPropagate(err error) error
	LogRequestInfo(info any)

	Info(strOrErr any)
	InfoPropagate(strOrErr any) error

	Debug(strOrErr any)
	DebugPropagate(strOrErr any) error

	Warning(strOrErr any)
	WarningPropagate(strOrErr any) error

	Error(strOrErr any)
	ErrorPropagate(strOrErr any) error

	Critical(strOrErr any)
	CriticalPropagate(strOrErr any) error

	Emergency(strOrErr any)
	EmergencyPropagate(strOrErr any) error

	SetOutput(w io.Writer)
	SetContext(ctx context.Context)

	Close() func()
}
