package errtype

const (
	info      = "INFO"
	debug     = "DEBUG"
	warning   = "WARNING"
	er        = "ERROR"
	critical  = "CRITICAL"
	emergency = "EMERGENCY"
)

type InfoLevel struct {
	Message  string
	File     string
	Function string
	Line     int
}

func (e InfoLevel) Error() string {
	return e.Message
}
func (e InfoLevel) Type() string {
	return info
}

type DebugLevel struct {
	Message  string
	File     string
	Function string
	Line     int
}

func (e DebugLevel) Error() string {
	return e.Message
}
func (e DebugLevel) Type() string {
	return debug
}

type WarningLevel struct {
	Message  string
	File     string
	Function string
	Line     int
}

func (e WarningLevel) Error() string {
	return e.Message
}
func (e WarningLevel) Type() string {
	return warning
}

type ErrorLevel struct {
	Err      error
	File     string
	Function string
	Line     int
}

func (e ErrorLevel) Error() string {
	return e.Err.Error()
}
func (e ErrorLevel) Type() string {
	return er
}

type CriticalLevel struct {
	Err      error
	File     string
	Function string
	Line     int
}

func (e CriticalLevel) Error() string {
	return e.Err.Error()
}
func (e CriticalLevel) Type() string {
	return critical
}

type EmergencyLevel struct {
	Err      error
	File     string
	Function string
	Line     int
}

func (e EmergencyLevel) Error() string {
	return e.Err.Error()
}
func (e EmergencyLevel) Type() string {
	return emergency
}
