package logger

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Warning(msg string)
	Error(err error)
	Critical(err error)
	Emergency(err error)
}
