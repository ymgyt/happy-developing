package app

import "context"

// Tag represents blog post tags.
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// TagToCreate -
type TagToCreate struct {
	Name string
}

// ListTagsInput -
type ListTagsInput struct {
	Name string
}

// TagService -
type TagService interface {
	Create(context.Context, *TagToCreate) (*Tag, error)
	List(context.Context, *ListTagsInput) ([]*Tag, error)
}
