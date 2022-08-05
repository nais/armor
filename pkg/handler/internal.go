package handler

import "net/http"

const (
	EndpointIsAlive = "/internal/isalive"
	EndpointIsReady = "/internal/isready"
)

func (h *Handler) isAlive(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("alive"))
}

func (h *Handler) isReady(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("ready"))
}
