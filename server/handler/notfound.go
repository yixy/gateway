package handler

import (
	"net/http"
	"time"
)

func NotFoundHandler(w http.ResponseWriter, req *http.Request) {
	time.Sleep(15 * time.Second)
	_, err := w.Write([]byte("404"))
	if err != nil {
	}
}
