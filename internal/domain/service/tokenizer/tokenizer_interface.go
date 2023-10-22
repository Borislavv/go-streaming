package tokenizer

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"net/http"
)

type Tokenizer interface {
	New(user *agg.User) (token string, err error)
	Has(r *http.Request) bool
	Get(r *http.Request) (token string, err error)
	Set(w http.ResponseWriter, r *http.Request, user *agg.User) error
	Refresh(w http.ResponseWriter, user *agg.User) error
	IsValid(w http.ResponseWriter, r *http.Request, user *agg.User) (ok bool, err error)
	Remove(w http.ResponseWriter)
}
