package serverhttp

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/server"
)

// middleware представляет структуру для обработки middleware с логгером и http.Handler.
type middleware struct {
	logger  server.Logger
	Handler http.Handler
}

// newMiddleware создает новый middleware с логгером и обработчиком HTTP-запросов.
func newMiddleware(logger server.Logger, httpHandler http.Handler) *middleware {
	return &middleware{
		logger:  logger,
		Handler: httpHandler,
	}
}

// logging добавляет middleware для логирования запросов и ответов.
func (m *middleware) logging() *middleware {
	curHandler := m.Handler

	m.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Создаем обертку для ResponseWriter, чтобы отслеживать статус-код.
		wrappedWriter := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		curHandler.ServeHTTP(wrappedWriter, r)

		handleTime := time.Since(start)

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			m.logger.Error(fmt.Sprintf("error splitting host and port, remote address: %s", r.RemoteAddr))
			ip = r.RemoteAddr
		}

		statusCode := wrappedWriter.statusCode

		logMessage := fmt.Sprintf("%s [%s] %s %s %s %d %s %s",
			ip, start.Format(time.RFC1123), r.Method, r.URL, r.Proto, statusCode, handleTime, r.UserAgent())
		m.logger.Info(logMessage)
	})

	return m
}

// responseWriter представляет обертку для http.ResponseWriter для отслеживания статус-кода.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader переопределяет метод WriteHeader для отслеживания статус-кода.
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
