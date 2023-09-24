package errors

// errored is a base error struct, implements sketch of necessary fields and base functionality.
type errored struct {
	ErrorMessage string `json:"message"` // Error is error interface implementation.
	ErrorType    string `json:"type"`    // Type of part of application: validation for example.
	errorStatus  int    // Status can be represented as http or stdout exit status and so on.
	errorLevel   int    // Level of an error, represents as logger.IOTA, see logger.errors file.
}

func (e errored) Error() string {
	return e.ErrorMessage
}
func (e errored) Status() int {
	return e.errorStatus
}
func (e errored) Level() int {
	return e.errorLevel
}
func (e errored) Type() string {
	return e.ErrorType
}

type publicError struct{ errored }

func (e publicError) Public() bool {
	return true
}

func IsPublic(err error) bool {
	_, ok := err.(publicError)
	return ok
}

type internalError struct{ errored }

func (e internalError) Internal() bool {
	return true
}

func IsInternal(err error) bool {
	_, ok := err.(internalError)
	return ok
}
