package security

type PasswordHasher struct {
}

func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{}
}

func (s *PasswordHasher) Hash(password string) (hash string, err error) {
	return "", nil
}

func (s *PasswordHasher) Verify(password string) (err error) {
	return nil
}
