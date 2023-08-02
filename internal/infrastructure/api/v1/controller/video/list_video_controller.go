package video

import (
	"log"
	"net/http"
)

const ListPath = "/video"

func List(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Hello world from GET:list method!")); err != nil {
		log.Fatalln(err)
	}
}
