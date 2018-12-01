package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ymgyt/happy-developing/hpdev/app"
)

// Static -
type Static struct {
	fs http.Handler

	*app.Env
	*base
}

// StaticRoot -
func (s *Static) StaticRoot(root string, prefix string) *Static {
	s.fs = http.StripPrefix(prefix, http.FileServer(http.Dir(root)))

	return s
}

// ServeStatic -
func (s *Static) ServeStatic(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s.fs.ServeHTTP(w, r)
}
