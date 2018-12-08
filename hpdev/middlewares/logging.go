package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/ymgyt/happy-developing/hpdev/app"
)

// Logging -
type Logging struct {
	Env     *app.Env
	logger  *zap.Logger
	sugar   *zap.SugaredLogger
	logging func(code int, r *http.Request, elapsed time.Duration)
	next    http.Handler
}

// MustLogging -
func MustLogging(env *app.Env) *Logging {
	l, err := NewLogging(env)
	if err != nil {
		panic(err)
	}
	return l
}

// NewLogging -
func NewLogging(env *app.Env) (*Logging, error) {
	l := &Logging{Env: env, logger: env.Log}
	l.logging = l.stdLogging

	if env.Mode == app.DevelopmentMode {
		// addCaller optionを適用したくない
		l.sugar = zap.New(l.logger.Core()).Sugar()
		l.logging = l.console

	}
	return l, nil
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
	m.logging(lw.statusCode, r, elapsed)
}

func (m *Logging) stdLogging(code int, r *http.Request, elapsed time.Duration) {
	if code >= 400 {
		// Errorで出すと、logic側のloggingとstack traceが重複する. optionの指定でなんとかできるかもしれない.
		m.logger.Warn("req",
			zap.Int("c", code),
			zap.String("m", r.Method),
			zap.String("u", r.URL.String()),
			zap.Float64("et", elapsed.Seconds()))
	} else {
		m.logger.Info("req",
			zap.Int("c", code),
			zap.String("m", r.Method),
			zap.String("u", r.URL.String()),
			zap.Float64("et", elapsed.Seconds()))
	}
}

func (m *Logging) console(code int, r *http.Request, elapsed time.Duration) {
	if strings.HasPrefix(r.URL.Path, "/static") {
		return
	}
	msg := fmt.Sprintf("|%3d| %-4s %-40s %.3f", code, r.Method, r.URL.String(), elapsed.Seconds())
	if code >= 400 {
		m.sugar.Warn(msg)
	} else {
		m.sugar.Info(msg)
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
