package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ymgyt/happy-developing/hpdev/app"
)

// Tag -
type Tag struct {
	*base
	service app.TagService
}

// Create -
func (t *Tag) Create(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	toCreate, err := t.readJSON(r)
	if err != nil {
		t.json(&apiResponse{W: w, Err: err})
		return
	}
	tag, err := t.service.Create(t.Env.Ctx, &app.TagToCreate{
		Name: toCreate.Name,
	})
	t.json(&apiResponse{W: w, Data: tag, Err: err})
}

// List -
func (t *Tag) List(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	tags, err := t.service.List(t.Env.Ctx, &app.ListTagsInput{})
	t.json(&apiResponse{W: w, Data: tags, Err: err})
}

func (t *Tag) readJSON(r *http.Request) (*app.Tag, error) {
	var tag app.Tag
	return &tag, json.NewDecoder(r.Body).Decode(&tag)
}
