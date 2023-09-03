package errs

type FieldCannotBeEmptyError struct {
	Message string
}

func NewFieldCannotBeEmptyError(msg string) FieldCannotBeEmptyError {
	return FieldCannotBeEmptyError{Message: msg}
}

func (e FieldCannotBeEmptyError) Error() string {
	return e.Message
}
