package middlewares

import (
	"net/http"
)

// Chain -
type Chain struct {
	middlewares []http.Handler
	router      http.Handler
}

// Middleware -
type Middleware interface {
	SetNext(http.Handler)
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// NewChain -
func NewChain(router http.Handler, mws ...Middleware) http.Handler {

	n := len(mws)
	for i := 0; i < n-1; i++ {
		mws[i].SetNext(mws[i+1])

	}

	mws[n-1].SetNext(router)

	return mws[0]
}
