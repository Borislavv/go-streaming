package dto_interface

type AuthRequest interface {
	GetEmail() string
	GetPassword() string
}
