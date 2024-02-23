package dtointerface

type AuthRequest interface {
	GetEmail() string
	GetPassword() string
}
