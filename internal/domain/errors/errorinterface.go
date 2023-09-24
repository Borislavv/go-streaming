package errors

// PublicError is the interface tells the responder that you want to show an error.
// All errors which implements it will show them messages and types.
type PublicError interface {
	Error() string
	Status() int
	Type() string
	Public() bool
}

type InternalError interface {
	Error() string
	Internal() bool
}
