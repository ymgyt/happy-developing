package gcp

import (
	"context"
	"time"

	"github.com/davecgh/go-spew/spew"

	"cloud.google.com/go/datastore"

	"github.com/ymgyt/happy-developing/hpdev/app"
)

const (
	postKind      = "post"
	postCreatedAt = "CreatedAt"
)

// PostStore implements app.PostService.
type PostStore struct {
	env *app.Env
	ds  *datastore.Client
}

// NewPostStore -
func NewPostStore(env *app.Env, ds *datastore.Client) (*PostStore, error) {
	return &PostStore{env: env, ds: ds}, nil
}

type post struct {
	ID           string
	Title        string
	URLSafeTitle string
	Content      string
	CreatedAt    time.Time
	Tags         []app.Tag
}

// Create -
func (store *PostStore) Create(ctx context.Context, post *app.Post) (*app.Post, error) {
	key := datastore.IncompleteKey(postKind, nil)

	newKey, err := store.ds.Put(ctx, key, store.toDatastoreRepl(post))
	// TODO Wrap error
	if err != nil {
		return nil, err
	}
	post.ID = newKey.Encode()

	return post, nil
}

// ListPosts -
func (store *PostStore) ListPosts(ctx context.Context, input *app.ListPostsInput) ([]*app.Post, error) {
	query := datastore.NewQuery(postKind).Offset(input.Offset).Limit(input.Limit)
	switch input.Order.Type {
	case app.PostOrderCreated:
		query = query.Order(postCreatedAt)
	default:
		query = query.Order(postCreatedAt)
	}

	var posts []*post
	keys, err := store.ds.GetAll(ctx, query, &posts)
	if err != nil {
		return nil, err
	}

	spew.Dump("got posts", posts)
	spew.Dump("got keys", keys)

	return store.fromDatastoreRepls(posts), nil
}

func (store *PostStore) toDatastoreRepl(p *app.Post) *post {
	return &post{
		ID:           p.ID,
		Title:        p.Title,
		URLSafeTitle: p.URLSafeTitle,
		Content:      p.Content,
		CreatedAt:    p.CreatedAt,
		Tags:         p.Tags,
	}
}

func (store *PostStore) fromDatastoreRepl(r *post) *app.Post {
	return &app.Post{
		ID:           r.ID,
		Title:        r.Title,
		URLSafeTitle: r.URLSafeTitle,
		Content:      r.Content,
		CreatedAt:    r.CreatedAt,
		Tags:         r.Tags,
	}
}

func (store *PostStore) fromDatastoreRepls(repls []*post) []*app.Post {
	var posts = make([]*app.Post, len(repls))
	for i := range repls {
		posts[i] = store.fromDatastoreRepl(repls[i])
	}
	return posts
}
