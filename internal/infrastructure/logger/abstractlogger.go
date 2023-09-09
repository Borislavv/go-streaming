package logger

import (
	"errors"
	"runtime"
)

type abstractLogger struct {
}

func (l *abstractLogger) trace() (file string, function string, line int) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	return frame.File, frame.Func.Name(), frame.Line
}

func (l *abstractLogger) error(strOrErr any) error {
	err, isErr := strOrErr.(error)
	if isErr {
		return err
	} else {
		str, isStr := strOrErr.(string)
		if isStr {
			return errors.New(str)
		}
	}
	panic("abstractlogger.error(): logging data is not a string or error type")
}

func (l *abstractLogger) toReadableLevel(err introspectedError) string {
	switch err.Level() {
	case InfoLevel:
		return InfoLevelReadable
	case DebugLevel:
		return DebugLevelReadable
	case WarningLevel:
		return WarningLevelReadable
	case ErrorLevel:
		return ErrorLevelReadable
	case CriticalLevel:
		return CriticalLevelReadable
	case EmergencyLevel:
		return EmergencyLevelReadable
	}
	panic("abstractlogger.toReadableLevel(): received undefined error level")
}

func (l *abstractLogger) toLevel(readableLevel string) int {
	switch readableLevel {
	case InfoLevelReadable:
		return InfoLevel
	case DebugLevelReadable:
		return DebugLevel
	case WarningLevelReadable:
		return WarningLevel
	case ErrorLevelReadable:
		return ErrorLevel
	case CriticalLevelReadable:
		return CriticalLevel
	case EmergencyLevelReadable:
		return EmergencyLevel
	}
	panic("abstractlogger.toLevel(): received undefined readable level")
}
