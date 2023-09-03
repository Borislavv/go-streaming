package errs

type NotFoundError struct {
	Message string
}

func NewNotFoundError(msg string) NotFoundError {
	return NotFoundError{Message: msg}
}

func (e NotFoundError) Error() string {
	return e.Message
}
