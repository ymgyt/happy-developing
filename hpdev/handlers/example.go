package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/ymgyt/happy-developing/hpdev/app"
)

// Example -
type Example struct {
	ts           *templateSet
	templateName string

	*app.Env
	*base
}

// RenderExample -
func (e *Example) RenderExample(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := e.ts.ExecuteTemplate(w, e.templateName, nil)
	e.handleRenderError(err)
}
