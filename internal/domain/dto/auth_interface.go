package dto

type AuthRequest interface {
	GetEmail() string
	GetPassword() string
}
