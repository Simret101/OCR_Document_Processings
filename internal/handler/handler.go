package handler

import "net/http"

type Upload interface {
	UploadStatmentHandler(w http.ResponseWriter, r *http.Request)
}
