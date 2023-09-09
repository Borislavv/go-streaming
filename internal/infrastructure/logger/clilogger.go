package logger

import (
	"log"
)

type CliLogger struct {
	abstractLogger
	errCh chan introspectedError
}

func NewCliLogger(buffer int) (logger *CliLogger, close func()) {
	l := &CliLogger{
		errCh: make(chan introspectedError, buffer),
	}
	l.handle()
	return l, l.Close()
}

func (l *CliLogger) Close() func() {
	return func() {
		close(l.errCh)
	}
}

func (l *CliLogger) Log(err error) {
	file, function, line := l.trace()
	l.log(err, file, function, line)
}

func (l *CliLogger) LogPropagate(err error) error {
	file, function, line := l.trace()
	l.log(err, file, function, line)
	return err
}

func (l *CliLogger) Info(strOrErr any) {
	file, function, line := l.trace()

	err := l.error(strOrErr)

	l.errCh <- infoLevelError{
		Er: err,
		Fl: file,
		Fn: function,
		Ln: line,
	}
}

func (l *CliLogger) InfoPropagate(strOrErr any) error {
	file, function, line := l.trace()

	err := l.error(strOrErr)

	l.errCh <- infoLevelError{
		Er: err,
		Fl: file,
		Fn: function,
		Ln: line,
	}

	return err
}

func (l *CliLogger) Debug(strOrErr any) {
	file, function, line := l.trace()

	err := l.error(strOrErr)

	l.errCh <- debugLevelError{
		Er: err,
		Fl: file,
		Fn: function,
		Ln: line,
	}
}

func (l *CliLogger) DebugPropagate(strOrErr any) error {
	file, function, line := l.trace()

	err := l.error(strOrErr)

	l.errCh <- debugLevelError{
		Er: err,
		Fl: file,
		Fn: function,
		Ln: line,
	}

	return err
}

func (l *CliLogger) Warning(strOrErr any) {
	file, function, line := l.trace()

	err := l.error(strOrErr)

	l.errCh <- warningLevelError{
		Er: err,
		Fl: file,
		Fn: function,
		Ln: line,
	}
}

func (l *CliLogger) WarningPropagate(strOrErr any) error {
	file, function, line := l.trace()

	err := l.error(strOrErr)

	l.errCh <- warningLevelError{
		Er: err,
		Fl: file,
		Fn: function,
		Ln: line,
	}

	return err
}

func (l *CliLogger) Error(strOrErr any) {
	file, function, line := l.trace()

	err := l.error(strOrErr)

	l.errCh <- errorLevelError{
		Er: err,
		Fl: file,
		Fn: function,
		Ln: line,
	}
}

func (l *CliLogger) ErrorPropagate(strOrErr any) error {
	file, function, line := l.trace()

	err := l.error(strOrErr)

	l.errCh <- errorLevelError{
		Er: err,
		Fl: file,
		Fn: function,
		Ln: line,
	}

	return err
}

func (l *CliLogger) Critical(strOrErr any) {
	file, function, line := l.trace()

	err := l.error(strOrErr)

	l.errCh <- criticalLevelError{
		Er: err,
		Fl: file,
		Fn: function,
		Ln: line,
	}
}

func (l *CliLogger) CriticalPropagate(strOrErr any) error {
	file, function, line := l.trace()

	err := l.error(strOrErr)

	l.errCh <- criticalLevelError{
		Er: err,
		Fl: file,
		Fn: function,
		Ln: line,
	}

	return err
}

func (l *CliLogger) Emergency(strOrErr any) {
	file, function, line := l.trace()

	err := l.error(strOrErr)

	l.errCh <- emergencyLevelError{
		Er: err,
		Fl: file,
		Fn: function,
		Ln: line,
	}
}

func (l *CliLogger) EmergencyPropagate(strOrErr any) error {
	file, function, line := l.trace()

	err := l.error(strOrErr)

	l.errCh <- emergencyLevelError{
		Er: err,
		Fl: file,
		Fn: function,
		Ln: line,
	}

	return err
}

// TODO must be injected writer for write logs through it
func (l *CliLogger) handle() {
	go func() {
		for err := range l.errCh {
			log.Printf(
				"\n\t[%v]\n\t\tFile: %v:%d\n\t\tFunc: %v\n\t\tMessage: %v\n",
				l.toReadableLevel(err),
				err.File(),
				err.Line(),
				err.Func(),
				err.Error(),
			)
		}
	}()
}

func (l *CliLogger) log(e error, file string, function string, line int) {
	err, isLoggableErr := e.(LoggableError)
	if !isLoggableErr {
		l.errCh <- errorLevelError{
			Er: e,
			Fl: file,
			Fn: function,
			Ln: line,
		}
		return
	}

	switch err.Level() {
	case InfoLevel:
		l.errCh <- infoLevelError{
			Er: err,
			Fl: file,
			Fn: function,
			Ln: line,
		}
		return
	case DebugLevel:
		l.errCh <- infoLevelError{
			Er: err,
			Fl: file,
			Fn: function,
			Ln: line,
		}
		return
	case WarningLevel:
		l.errCh <- warningLevelError{
			Er: err,
			Fl: file,
			Fn: function,
			Ln: line,
		}
		return
	case ErrorLevel:
		l.errCh <- errorLevelError{
			Er: err,
			Fl: file,
			Fn: function,
			Ln: line,
		}
		return
	case CriticalLevel:
		l.errCh <- criticalLevelError{
			Er: err,
			Fl: file,
			Fn: function,
			Ln: line,
		}
		return
	case EmergencyLevel:
		l.errCh <- emergencyLevelError{
			Er: err,
			Fl: file,
			Fn: function,
			Ln: line,
		}
		return
	}

	panic("clilogger.log(): undefined error level received")
}
