package controller

import (
	"github.com/gorilla/mux"
)

type Controller interface {
	// AddRoute is a method which must take a *mux.Router and add new route.
	// Must be implemented into each RestApi and so on controllers for inject them
	// into the server for serving them.
	//
	// Commonly may be represented as:
	// 	router.
	//		Path(CreatePath).
	//		HandlerFunc(c.Create).
	//		Methods(http.MethodPost)
	AddRoute(router *mux.Router)
}
