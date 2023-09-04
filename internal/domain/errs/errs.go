package errs

type PublicError interface {
	Error() string
	Status() int
}
