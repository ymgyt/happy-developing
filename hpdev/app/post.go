package app

import (
	"context"
	"time"
)

// Post represents Blog Post.
type Post struct {
	ID           string
	Title        string
	URLSafeTitle string // supposed to be users url. e.g /posts/safe_title.
	Content      string
	CreatedAt    time.Time
	Tags         []Tag
}

// PostOrderType -
type PostOrderType string

const (
	// PostOrderCreated -
	PostOrderCreated PostOrderType = "created"
)

// PostOrder -
type PostOrder struct {
	Type PostOrderType
	Desc bool
}

// ListPostsInput -
type ListPostsInput struct {
	Limit  int
	Offset int
	Order  PostOrder
}

// PostService is responsible for post data CRUD.
type PostService interface {
	Create(context.Context, *Post) (*Post, error)
	ListPosts(context.Context, *ListPostsInput) ([]*Post, error)
}
