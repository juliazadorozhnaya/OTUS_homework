package handler

import (
	"fmt"
	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/server/server"
	"io"
	"net/http"
)

type Handler struct {
	logger internalhttp.Logger
}

func NewHandler(logger internalhttp.Logger) *Handler {
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
