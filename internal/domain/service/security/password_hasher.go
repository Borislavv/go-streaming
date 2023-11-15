package security

type PasswordHasher interface {
	Hash(password string) (hash string, err error)
	Verify(password string) (err error)
}