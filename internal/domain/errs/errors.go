package errs

// PublicError is the interface tells the responder that you want to show an error.
// All errors which implements it will show them messages and types
type PublicError interface {
	Error() string
	Status() int
	Type() string
}

// errored is a base error struct, implements sketch of necessary fields and base functionality
type errored struct {
	ErrorMessage string `json:"message"` // Error is error interface implementation.
	ErrorType    string `json:"type"`    // Type of part of application: validation for example.
	errorStatus  int    // Status can be represented as http or cli exit status and so on.
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
