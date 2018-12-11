package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Router -
type Router struct {
	tree     map[string]*httprouter.Router
	mws      []http.Handler
	children map[string]*Router
}

// Group -
func (r *Router) Group(prefix string) *Router {
	if r.children == nil {
		r.children = make(map[string]*Router)
	}
	child := r.children[prefix]
	if child != nil {
		panic("group conflict")
	}
	// middlewareを引き継ぐ
	child = &Router{mws: r.mws}
	r.children[prefix] = child

	return child
}

// GET -
func (r *Router) GET(path string, h httprouter.Handle) {
	r.addHandler("GET", path, h)
}

// POST -
func (r *Router) POST(path string, h httprouter.Handle) {
	r.addHandler("POST", path, h)
}

// ServeHTTP -
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h := r.tree[req.Method]
	if h == nil {
		r.MethodNotAllowed(w, req)
	}
	// apply middlewares
	for _, mw := range r.mws {
		mw.ServeHTTP(w, req)
	}
	h.ServeHTTP(w, req)
}

// MethodNotAllowed -
func (r *Router) MethodNotAllowed(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (r *Router) addHandler(method, path string, h httprouter.Handle) {
	if r.tree == nil {
		r.tree = make(map[string]*httprouter.Router)
	}

	hr := r.tree[method]
	if hr == nil {
		hr = httprouter.New()
		r.tree[method] = hr
	}

	hr.Handle(method, path, h)
}
