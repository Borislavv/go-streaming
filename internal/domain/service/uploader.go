package service

import "net/http"

type Uploader interface {
	Upload(r *http.Request) (resourceId string, err error)
}
