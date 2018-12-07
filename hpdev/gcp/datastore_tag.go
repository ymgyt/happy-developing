package gcp

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/ymgyt/happy-developing/hpdev/app"
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
	// TODO check already exist tags
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

// List -
func (store *TagStore) List(ctx context.Context, in *app.ListTagsInput) ([]*app.Tag, error) {
	q := datastore.NewQuery(tagKind)

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
