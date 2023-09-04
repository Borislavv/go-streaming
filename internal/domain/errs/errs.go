package errs

type PublicError interface {
	Error() string
	Public() bool
}
