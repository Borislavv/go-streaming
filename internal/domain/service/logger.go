package service

// TODO - logger must be improved for write 3-5 or all last stack frames due to handle errors properly
type Logger interface {
	Info(msg string)
	Debug(msg string)
	Warning(msg string)
	Error(err error)
	Critical(err error)
	Emergency(err error)
}
