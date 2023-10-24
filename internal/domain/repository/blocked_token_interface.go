package repository

type BlockedToken interface {
	Has(token string) (found bool, err error)
}
