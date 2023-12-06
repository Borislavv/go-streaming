package _interface

type AuthRequest interface {
	GetEmail() string
	GetPassword() string
}
