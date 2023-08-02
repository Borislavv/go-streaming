package video

import (
	"log"
	"net/http"
)

const UpdatePath = "/video/{id}"

func Update(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Hello world from PATCH method!")); err != nil {
		log.Fatalln(err)
	}
}
