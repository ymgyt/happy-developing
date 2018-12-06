package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"

	"github.com/ymgyt/happy-developing/hpdev/app"
)

// Post はブログの記事を管理する責務をもつ.
type Post struct {
	templateName string
	ts           *templateSet
	*base

	service app.PostService
}

// RenderMetaList -
func (p *Post) RenderMetaList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	metaList, err := p.service.ListMeta(p.Env.Ctx, &app.ListMetaInput{
		Limit:  100,
		Offset: 0,
		Order:  app.PostOrder{Type: app.PostOrderCreated, Desc: true},
		Status: nil,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = p.ts.ExecuteTemplate(w, "author/posts", metaList)
	p.handleRenderError(err)
}

// RenderPostForm -
func (p *Post) RenderPostForm(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	err := p.ts.ExecuteTemplate(w, "author/post", nil)
	p.handleRenderError(err)
}

// Create -
func (p *Post) Create(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	post, err := p.readJSON(r)
	if err != nil {
		p.json(&apiResponse{W: w, Err: err})
		return
	}
	created, err := p.service.Create(p.Env.Ctx, post)
	p.json(&apiResponse{W: w, Data: created, Err: err})
}

// Get -
func (p *Post) Get(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	metaID, err := p.readMetaID(params)
	if err != nil {
		p.invalidRequest(w, err)
		return
	}

	post, err := p.service.Get(p.Env.Ctx, &app.GetPostInput{MetaID: metaID})
	p.json(&apiResponse{W: w, Data: post, Err: err})
}

// Update -
func (p *Post) Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	post, err := p.readJSON(r)
	if err != nil {
		panic(err)
	}
	_ = post

	w.Write([]byte("OK"))
}

func (p *Post) readMetaID(params httprouter.Params) (metaID int64, err error) {
	return strconv.ParseInt(params.ByName("metaid"), 10, 64)
}

func (p *Post) readPost(r *http.Request) (*app.Post, error) {
	rawMetaKeyID := r.PostFormValue("post-key-id")
	metaKeyID, err := strconv.ParseInt(rawMetaKeyID, 10, 64)
	if err != nil {
		p.Env.Log.Error("handle_post", zap.String("rawMetakeyID", rawMetaKeyID), zap.Error(err))
		return nil, err
	}
	meta := &app.PostMeta{
		Key:          &datastore.Key{ID: metaKeyID},
		Title:        r.PostFormValue("post-title"),
		URLSafeTitle: r.PostFormValue("post-url-safe-title"),
		CreatedAt:    p.Env.Now(),
	}

	content := &app.PostContent{
		HTML: r.PostFormValue("post-content"),
	}

	return &app.Post{Meta: meta, Content: content}, nil
}

func (p *Post) readJSON(r *http.Request) (*app.Post, error) {
	var post app.Post
	return &post, json.NewDecoder(r.Body).Decode(&post)
}
