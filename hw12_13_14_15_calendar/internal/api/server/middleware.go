package server

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

type Middleware struct {
	logger  Logger
	Handler http.Handler
}

func NewMiddleware(logger Logger, httpHandler http.Handler) *Middleware {
	return &Middleware{
		logger:  logger,
		Handler: httpHandler,
	}
}

func (m *Middleware) Logging() *Middleware {
	curHandler := m.Handler

	m.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		curHandler.ServeHTTP(w, r)
		handleTime := time.Since(start)

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			m.logger.Error(fmt.Sprintf("error split host port, remote address:%s", r.RemoteAddr))
		}

		statusCode := w.Header().Get("status")

		logMessage := fmt.Sprintf("%s [%s] %s %s %s %s %s %s",
			ip, start, r.Method, r.URL, r.Proto, statusCode, handleTime, r.UserAgent())
		m.logger.Info(logMessage)
	})

	return m
}
