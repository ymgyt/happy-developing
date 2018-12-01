package handlers

import (
	"net/http"

	"github.com/davecgh/go-spew/spew"

	"github.com/julienschmidt/httprouter"

	"github.com/ymgyt/happy-developing/hpdev/app"
)

// Post はブログの記事を管理する責務をもつ.
type Post struct {
	templateName string
	ts           *templateSet
	*base

	service app.PostService
}

// RenderPostForm -
func (p *Post) RenderPostForm(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	err := p.ts.ExecuteTemplate(w, p.templateName, nil)
	p.handleRenderError(err)
}

// CreatePost -
func (p *Post) CreatePost(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	post, err := p.readPost(r)
	if err != nil {
		// TODO organize api response handling
		http.Error(w, err.Error(), 500)
		return
	}

	created, err := p.service.Create(p.Env.Ctx, post)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	p.Env.Log.Debug("post created", created)
	w.Write([]byte("OK"))
}

// ListPosts -
func (p *Post) ListPosts(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	p.Env.Log.Debug("list posts")

	// TODO populate input param
	posts, err := p.service.ListPosts(p.Env.Ctx, &app.ListPostsInput{
		Limit:  100,
		Offset: 0,
		Order:  app.PostOrder{Type: app.PostOrderCreated, Desc: true},
	})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	spew.Fdump(w, posts)
}

func (p *Post) readPost(r *http.Request) (*app.Post, error) {
	post := &app.Post{
		Title:        r.PostFormValue("post-title"),
		URLSafeTitle: r.PostFormValue("post-url-safe-title"),
		Content:      r.PostFormValue("post-content"),
		CreatedAt:    p.Env.Now(),
		// TODO tags
	}

	// TODO validate here
	return post, nil
}
