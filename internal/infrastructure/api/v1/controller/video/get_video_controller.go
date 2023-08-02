package video

import (
	"log"
	"net/http"
)

const GetPath = "/video/{id}"

func Get(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Hello world from GET:item method!")); err != nil {
		log.Fatalln(err)
	}
}
