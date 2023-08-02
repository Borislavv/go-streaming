package video

import (
	"log"
	"net/http"
)

const CreatePath = "/video"

func Create(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Hello world from POST method!")); err != nil {
		log.Fatalln(err)
	}
}
