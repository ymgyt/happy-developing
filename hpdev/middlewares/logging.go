package middlewares

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Logging -
type Logging struct {
	next http.Handler
	Log  *zap.Logger
}

// ServeHTTP -
func (m *Logging) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lw := &loggingWriter{ResponseWriter: w}
	start := time.Now()
	m.next.ServeHTTP(lw, r)
	elapsed := time.Since(start)

	// この時点で書き込まれていない場合がある。
	if lw.statusCode == 0 {
		lw.statusCode = 200
	}

	if lw.statusCode >= 400 {
		// Errorで出すと、logic側のloggingとstack traceが重複する.
		m.Log.Warn("req",
			zap.Int("c", lw.statusCode),
			zap.String("m", r.Method),
			zap.String("u", r.URL.String()),
			zap.Float64("et", elapsed.Seconds()))
	} else {
		m.Log.Info("req",
			zap.Int("c", lw.statusCode),
			zap.String("m", r.Method),
			zap.String("u", r.URL.String()),
			zap.Float64("et", elapsed.Seconds()))
	}
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
