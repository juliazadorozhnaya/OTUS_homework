package server

import (
	"fmt"
	"io"
	"net/http"
)

type Handler struct {
	logger Logger
}

func NewHandler(logger Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func (h *Handler) GetHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("status", fmt.Sprint(http.StatusOK))
	w.WriteHeader(http.StatusOK)

	_, err := io.WriteString(w, "hello-world\n")
	if err != nil {
		h.logger.Error(err.Error())
		return
	}
}
