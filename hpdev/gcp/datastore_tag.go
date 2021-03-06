package gcp

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
	"github.com/ymgyt/happy-developing/hpdev/app"
	"github.com/ymgyt/happy-developing/hpdev/errors"
)

const (
	tagKind = "Tag"
)

// TagStore -
type TagStore struct {
	env *app.Env
	ds  *datastore.Client
}

type tag struct {
	*app.Tag
	Key *datastore.Key `datastore:"__key__"`
}

// NewTagStore -
func NewTagStore(env *app.Env, ds *datastore.Client) (*TagStore, error) {
	return &TagStore{env: env, ds: ds}, nil
}

// Create -
func (store *TagStore) Create(ctx context.Context, tag *app.TagToCreate) (*app.Tag, error) {
	// 既存check
	exists, err := store.exists(ctx, tag.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.AlreadyExists(fmt.Sprintf("tag: %s", tag.Name), nil)
	}

	incKey := datastore.IncompleteKey(tagKind, nil)
	key, err := store.ds.Put(ctx, incKey, tag)
	if err != nil {
		return nil, err
	}
	return &app.Tag{
		ID:   int(key.ID),
		Name: tag.Name,
	}, nil
}

func (store *TagStore) exists(ctx context.Context, name string) (bool, error) {
	tags, err := store.List(ctx, &app.ListTagsInput{Name: name})
	return len(tags) >= 1, err
}

// List -
func (store *TagStore) List(ctx context.Context, in *app.ListTagsInput) ([]*app.Tag, error) {
	q := datastore.NewQuery(tagKind)
	if in.Name != "" {
		q = q.Filter("Name =", in.Name)
	}

	var repls []*tag
	_, err := store.ds.GetAll(ctx, q, &repls)
	if err != nil {
		return nil, err
	}

	tags := make([]*app.Tag, len(repls))
	for i, r := range repls {
		tags[i] = r.toDomain()
	}
	return tags, nil
}

func (t *tag) toDomain() *app.Tag {
	t.Tag.ID = int(t.Key.ID)
	return t.Tag
}
