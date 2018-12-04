package gcp

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"

	"github.com/ymgyt/happy-developing/hpdev/app"
)

const (
	postMetaKind    = "PostMeta"
	postContentKind = "PostContent"
	postCreatedAt   = "CreatedAt"
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

// Create -
func (store *PostStore) Create(ctx context.Context, post *app.Post) (*app.Post, error) {
	// TODO transaction
	meta, err := store.createMeta(ctx, post.Meta)
	if err != nil {
		return nil, err
	}

	content, err := store.createContent(ctx, meta.Key.ID, post.Content)
	if err != nil {
		return nil, err
	}

	return &app.Post{Meta: meta, Content: content}, nil
}

func (store *PostStore) createMeta(ctx context.Context, meta *app.PostMeta) (*app.PostMeta, error) {
	k := datastore.IncompleteKey(postMetaKind, nil)
	nk, err := store.ds.Put(ctx, k, meta)
	if err != nil {
		return nil, err
	}
	meta.Key = nk

	return meta, nil
}

func (store *PostStore) createContent(ctx context.Context, metaKeyID int64, content *app.PostContent) (*app.PostContent, error) {
	// meta keyをparentに指定するか悩む
	k := datastore.IDKey(postContentKind, metaKeyID, nil)
	_, err := store.ds.Put(ctx, k, content)
	if err != nil {
		return nil, err
	}

	return content, nil
}

// Update -
func (store *PostStore) Update(ctx context.Context, post *app.Post) (*app.Post, error) {
	meta, err := store.updateMeta(ctx, post.Meta)
	if err != nil {
		return nil, datastoreErr(err)
	}
	content, err := store.updateContent(ctx, meta.Key.ID, post.Content)

	return &app.Post{Meta: meta, Content: content}, datastoreErr(err)
}

func (store *PostStore) updateMeta(ctx context.Context, meta *app.PostMeta) (*app.PostMeta, error) {
	k := datastore.IDKey(postMetaKind, meta.Key.ID, nil)
	_, err := store.ds.Put(ctx, k, meta)
	return meta, err
}

func (store *PostStore) updateContent(ctx context.Context, metaKeyID int64, content *app.PostContent) (*app.PostContent, error) {
	k := datastore.IDKey(postContentKind, metaKeyID, nil)
	_, err := store.ds.Put(ctx, k, content)
	return content, err
}

// ListMeta -
func (store *PostStore) ListMeta(ctx context.Context, input *app.ListMetaInput) ([]*app.PostMeta, error) {
	q := datastore.NewQuery(postMetaKind).Offset(input.Offset).Limit(input.Limit)

	var order string
	switch input.Order.Type {
	case app.PostOrderCreated:
		order = postCreatedAt
	default:
		order = postCreatedAt
	}
	if input.Order.Desc {
		order = "-" + order
	}
	q = q.Order(order)

	var meta []*app.PostMeta
	_, err := store.ds.GetAll(ctx, q, &meta)

	return meta, err
}

// Get -
func (store *PostStore) Get(ctx context.Context, input *app.GetPostInput) (*app.Post, error) {
	meta, err := store.getMetaByID(ctx, input.MetaID)
	if err != nil {
		return nil, datastoreErr(err)
	}
	content, err := store.getContentByMetaID(ctx, meta.Key.ID)
	if err != nil {
		return nil, datastoreErr(err)
	}

	return &app.Post{Meta: meta, Content: content}, nil
}

func (store *PostStore) getMetaByID(ctx context.Context, metaID int64) (*app.PostMeta, error) {
	k := datastore.IDKey(postMetaKind, metaID, nil)
	var meta app.PostMeta
	return &meta, store.ds.Get(ctx, k, &meta)
}

func (store *PostStore) getContentByMetaID(ctx context.Context, metaID int64) (*app.PostContent, error) {
	k := datastore.IDKey(postContentKind, metaID, nil)
	var content app.PostContent
	return &content, store.ds.Get(ctx, k, &content)
}

func datastoreErr(err error) error {
	if err == datastore.ErrNoSuchEntity {
		// wrap error
		fmt.Println("not found error")
	}
	return err
}
