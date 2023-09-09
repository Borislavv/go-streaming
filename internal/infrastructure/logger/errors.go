package logger

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
	Error() string
	File() string
	Func() string
	Line() int
	Level() int
}

type introspectionError struct {
	Er error  // Error
	Fl string // File
	Fn string // Function
	Ln int    // Line
}

type infoLevelError introspectionError

func (e infoLevelError) Error() string {
	return e.Er.Error()
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

func (e debugLevelError) Error() string {
	return e.Er.Error()
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

func (e warningLevelError) Error() string {
	return e.Er.Error()
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

func (e errorLevelError) Error() string {
	return e.Er.Error()
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

func (e criticalLevelError) Error() string {
	return e.Er.Error()
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

func (e emergencyLevelError) Error() string {
	return e.Er.Error()
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
