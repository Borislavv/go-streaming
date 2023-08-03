package controller

import (
	"github.com/gorilla/mux"
)

type Controller interface {
	AddRoute(router *mux.Router)
}
