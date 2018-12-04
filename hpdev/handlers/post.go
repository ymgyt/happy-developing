package handlers

import (
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"cloud.google.com/go/datastore"
	"github.com/julienschmidt/httprouter"

	"github.com/ymgyt/happy-developing/hpdev/app"
	"github.com/ymgyt/happy-developing/hpdev/view"
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
// 新規投稿の場合と既存の編集の場合がある。
func (p *Post) RenderPostForm(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	rawMetaID := params.ByName("metaid")

	p.Env.Log.Debug("render_post_form", zap.String("metaid", rawMetaID))
	// 新規投稿
	if rawMetaID == "new" {
		// nilだとtemplate側のhandlingが面倒なので、空を渡しておく.
		err := p.ts.ExecuteTemplate(w, "author/post", view.AuthorPost{Data: initialPost, FormAction: "new"})
		p.handleRenderError(err)
		return
	}

	// 既存の編集
	metaID, err := strconv.ParseInt(rawMetaID, 10, 64)
	if err != nil {
		p.handleRenderError(err)
		return
	}
	post, err := p.service.Get(p.Env.Ctx, &app.GetPostInput{MetaID: metaID})
	if err != nil {
		p.handleRenderError(err)
		return
	}

	err = p.ts.ExecuteTemplate(w, "author/post", view.AuthorPost{Data: post, FormAction: rawMetaID})
	p.handleRenderError(err)
}

// SavePost -
func (p *Post) SavePost(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	rawMetaID := params.ByName("metaid")
	post, err := p.readPost(r)
	if err != nil {
		// TODO organize api response handling
		http.Error(w, err.Error(), 500)
		return
	}

	// 新規投稿
	if rawMetaID == "new" {
		created, err := p.service.Create(p.Env.Ctx, post)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		metaID := strconv.FormatInt(created.Meta.Key.ID, 10)
		p.Env.Log.Debug("post_created", zap.String("meta_id", metaID), zap.String("title", created.Meta.Title))
		// redirectしたほうがURLがきれいになる.
		err = p.ts.ExecuteTemplate(w, "author/post", view.AuthorPost{Data: post, FormAction: metaID})
		p.handleRenderError(err)
		return
	}

	// 既存のupdate
	post, err = p.service.Update(p.Env.Ctx, post)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	metaID := strconv.FormatInt(post.Meta.Key.ID, 10)
	p.Env.Log.Debug("post_updated", zap.String("meta_id", metaID), zap.String("title", post.Meta.Title))
	err = p.ts.ExecuteTemplate(w, "author/post", view.AuthorPost{Data: post, FormAction: metaID})
	p.handleRenderError(err)
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
		HTML: []byte(r.PostFormValue("post-content")),
	}

	return &app.Post{Meta: meta, Content: content}, nil
}

// 新規投稿の際にtemplate側に渡す用のPost
var initialPost = &app.Post{
	Meta: &app.PostMeta{
		Key: &datastore.Key{},
	},
	Content: &app.PostContent{},
}
