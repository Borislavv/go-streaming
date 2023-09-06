package errs

import "net/http"

const (
	DefaultErrorMessage    = "internal server error"
	DefaultErrorType       = "application"
	DefaultErrorStatusCode = http.StatusInternalServerError
)

type PublicError interface {
	Error() string
	Status() int
}

type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	status  int
}

func NewError(Message string, Type string, Status int) *Error {
	if Message == "" {
		Message = DefaultErrorMessage
	}
	if Type == "" {
		Type = DefaultErrorType
	}
	if Status == 0 {
		Status = DefaultErrorStatusCode
	}
	return &Error{
		Message: Message,
		Type:    Type,
		status:  Status,
	}
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Status() int {
	return e.status
}
