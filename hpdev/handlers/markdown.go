package handlers

import (
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ymgyt/happy-developing/hpdev/app"
)

// Markdown -
type Markdown struct {
	*app.Env
	*base
}

// ConvertHTML -
func (m *Markdown) ConvertHTML(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// TODO limit read size
	md, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	converter := &app.Markdown{}
	htm := converter.ConvertHTML(md)

	syntaxed, err := app.SyntaxHighlight(htm)
	if err != nil {
		panic(err)
	}

	w.Write(syntaxed)
}
