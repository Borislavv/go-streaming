package logger

import "time"

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
}

type introspectionError struct {
	Dt time.Time `json:"date"`
	Mg string    `json:"message"`
	Fl string    `json:"file"`
	Fn string    `json:"function"`
	Ln int       `json:"line"`
}

type infoLevelError introspectionError

func (e infoLevelError) Date() time.Time {
	return e.Dt
}
func (e infoLevelError) Error() string {
	return e.Mg
}
func (e infoLevelError) File() string {
	return e.Fl
}
func (e infoLevelError) Func() string {
	return e.Fn
}
func (e infoLevelError) Line() int {
	return e.Ln
}
func (e infoLevelError) Level() int {
	return InfoLevel
}

type debugLevelError introspectionError

func (e debugLevelError) Date() time.Time {
	return e.Dt
}
func (e debugLevelError) Error() string {
	return e.Mg
}
func (e debugLevelError) File() string {
	return e.Fl
}
func (e debugLevelError) Func() string {
	return e.Fn
}
func (e debugLevelError) Line() int {
	return e.Ln
}
func (e debugLevelError) Level() int {
	return DebugLevel
}

type warningLevelError introspectionError

func (e warningLevelError) Date() time.Time {
	return e.Dt
}
func (e warningLevelError) Error() string {
	return e.Mg
}
func (e warningLevelError) File() string {
	return e.Fl
}
func (e warningLevelError) Func() string {
	return e.Fn
}
func (e warningLevelError) Line() int {
	return e.Ln
}
func (e warningLevelError) Level() int {
	return WarningLevel
}

type errorLevelError introspectionError

func (e errorLevelError) Date() time.Time {
	return e.Dt
}
func (e errorLevelError) Error() string {
	return e.Mg
}
func (e errorLevelError) File() string {
	return e.Fl
}
func (e errorLevelError) Func() string {
	return e.Fn
}
func (e errorLevelError) Line() int {
	return e.Ln
}
func (e errorLevelError) Level() int {
	return ErrorLevel
}

type criticalLevelError introspectionError

func (e criticalLevelError) Date() time.Time {
	return e.Dt
}
func (e criticalLevelError) Error() string {
	return e.Mg
}
func (e criticalLevelError) File() string {
	return e.Fl
}
func (e criticalLevelError) Func() string {
	return e.Fn
}
func (e criticalLevelError) Line() int {
	return e.Ln
}
func (e criticalLevelError) Level() int {
	return CriticalLevel
}

type emergencyLevelError introspectionError

func (e emergencyLevelError) Date() time.Time {
	return e.Dt
}
func (e emergencyLevelError) Error() string {
	return e.Mg
}
func (e emergencyLevelError) File() string {
	return e.Fl
}
func (e emergencyLevelError) Func() string {
	return e.Fn
}
func (e emergencyLevelError) Line() int {
	return e.Ln
}
func (e emergencyLevelError) Level() int {
	return EmergencyLevel
}
