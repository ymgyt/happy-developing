package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// HealthCheck -
type HealthCheck struct {
}

// Beat -
func (hc *HealthCheck) Beat(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("OK"))
}
