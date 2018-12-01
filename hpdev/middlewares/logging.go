package middlewares

import (
	"fmt"
	"net/http"
)

// Logging -
type Logging struct {
	next http.Handler
}

// ServeHTTP -
func (m *Logging) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lw := &loggingWriter{ResponseWriter: w}
	m.next.ServeHTTP(lw, r)

	// この時点で書き込まれていない場合がある。
	if lw.statusCode == 0 {
		lw.statusCode = 200
	}
	fmt.Printf("[%3d] %5s %s\n", lw.statusCode, r.Method, r.URL.String())
}

// SetNext -
func (m *Logging) SetNext(next http.Handler) {
	m.next = next
}

type loggingWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lw *loggingWriter) WriteHeader(statusCode int) {
	lw.statusCode = statusCode
	lw.ResponseWriter.WriteHeader(statusCode)
}
