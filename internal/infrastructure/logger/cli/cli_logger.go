package cli

import (
	"github.com/Borislavv/video-streaming/internal/infrastructure/logger/errtype"
	"log"
	"runtime"
)

type Logger struct {
	errCh chan error
}

func NewLogger(errCh chan error) *Logger {
	l := &Logger{
		errCh: errCh,
	}

	l.handle()

	return l
}

func (l *Logger) Info(msg string) {
	file, function, line := l.trace()
	l.errCh <- errtype.InfoLevel{
		Message:  msg,
		File:     file,
		Function: function,
		Line:     line,
	}
}

func (l *Logger) Debug(msg string) {
	file, function, line := l.trace()
	l.errCh <- errtype.DebugLevel{
		Message:  msg,
		File:     file,
		Function: function,
		Line:     line,
	}
}

func (l *Logger) Warning(msg string) {
	file, function, line := l.trace()
	l.errCh <- errtype.WarningLevel{
		Message:  msg,
		File:     file,
		Function: function,
		Line:     line,
	}
}

func (l *Logger) Error(e error) {
	file, function, line := l.trace()
	l.errCh <- errtype.ErrorLevel{
		Err:      e,
		File:     file,
		Function: function,
		Line:     line,
	}
}

func (l *Logger) Critical(e error) {
	file, function, line := l.trace()
	l.errCh <- errtype.CriticalLevel{
		Err:      e,
		File:     file,
		Function: function,
		Line:     line,
	}
}

func (l *Logger) Emergency(e error) {
	file, function, line := l.trace()
	l.errCh <- errtype.EmergencyLevel{
		Err:      e,
		File:     file,
		Function: function,
		Line:     line,
	}
}

func (l *Logger) handle() {
	go func() {
		for er := range l.errCh {
			t, f, fn, ln := l.details(er)
			log.Printf("\n\t[%s]\n\t\tFile: %s:%d\n\t\tFunc: %s\n\t\tMessage: %s\n", t, f, ln, fn, er)
		}
	}()
}

func (l *Logger) trace() (string, string, int) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	return frame.File, frame.Func.Name(), frame.Line
}

func (l *Logger) details(e error) (t string, f string, fn string, ln int) {
	if iErr, iOk := e.(errtype.InfoLevel); iOk {
		return iErr.Type(), iErr.File, iErr.Function, iErr.Line
	}

	if dErr, dOk := e.(errtype.DebugLevel); dOk {
		return dErr.Type(), dErr.File, dErr.Function, dErr.Line
	}

	if wErr, wOk := e.(errtype.WarningLevel); wOk {
		return wErr.Type(), wErr.File, wErr.Function, wErr.Line
	}

	if eErr, eOk := e.(errtype.ErrorLevel); eOk {
		return eErr.Type(), eErr.File, eErr.Function, eErr.Line
	}

	if cErr, cOk := e.(errtype.CriticalLevel); cOk {
		return cErr.Type(), cErr.File, cErr.Function, cErr.Line
	}

	if emErr, emOk := e.(errtype.EmergencyLevel); emOk {
		return emErr.Type(), emErr.File, emErr.Function, emErr.Line
	}

	return "unknown error type", "", "", 0
}