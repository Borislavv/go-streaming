package logger_interface

import (
	"context"
	"io"
)

type Logger interface {
	Log(err error)
	LogPropagate(err error) error
	LogData(data any)

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

	Writer() io.Writer
	SetOutput(w io.Writer)
	Context() context.Context
	SetContext(ctx context.Context)

	Close() func()
}
