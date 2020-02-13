package handler

import (
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, req *http.Request) {
	_, err := w.Write([]byte("404"))
	if err != nil {
	}
}
