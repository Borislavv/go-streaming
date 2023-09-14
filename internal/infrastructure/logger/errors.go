package logger

import (
	"time"
)

const (
	InfoLevelReadable      = "INFO"
	DebugLevelReadable     = "DEBUG"
	WarningLevelReadable   = "WARNING"
	ErrorLevelReadable     = "ERROR"
	CriticalLevelReadable  = "CRITICAL"
	EmergencyLevelReadable = "EMERGENCY"
	InfoLevel              = iota
	DebugLevel
	WarningLevel
	ErrorLevel
	CriticalLevel
	EmergencyLevel
)

type LoggableError interface {
	Error() string
	Level() int
}

type introspectedError interface {
	Date() time.Time
	Error() string
	File() string
	Func() string
	Line() int
	Level() int
	RequestId() string
	SetRequestId(id string)
}

type introspectionError struct {
	Dt time.Time `json:"date"`
	Rq string    `json:"requestID,omitempty"`
	Mg string    `json:"message"`
	Fl string    `json:"file"`
	Fn string    `json:"function"`
	Ln int       `json:"line"`
}

func (e *introspectionError) Date() time.Time {
	return e.Dt
}
func (e *introspectionError) Error() string {
	return e.Mg
}
func (e *introspectionError) File() string {
	return e.Fl
}
func (e *introspectionError) Func() string {
	return e.Fn
}
func (e *introspectionError) Line() int {
	return e.Ln
}
func (e *introspectionError) Level() int {
	return ErrorLevel
}
func (e *introspectionError) RequestId() string {
	return e.Rq
}
func (e *introspectionError) SetRequestId(id string) {
	e.Rq = id
}

type infoLevelError struct{ introspectionError }

func (e *infoLevelError) Level() int {
	return InfoLevel
}

type debugLevelError struct{ introspectionError }

func (e *debugLevelError) Level() int {
	return DebugLevel
}

type warningLevelError struct{ introspectionError }

func (e *warningLevelError) Level() int {
	return WarningLevel
}

type errorLevelError struct{ introspectionError }

func (e *errorLevelError) Level() int {
	return ErrorLevel
}

type criticalLevelError struct{ introspectionError }

func (e *criticalLevelError) Level() int {
	return CriticalLevel
}

type emergencyLevelError struct{ introspectionError }

func (e *emergencyLevelError) Level() int {
	return EmergencyLevel
}
