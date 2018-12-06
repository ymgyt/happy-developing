package app

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
)

// Post -
type Post struct {
	Meta    *PostMeta    `json:"meta"`
	Content *PostContent `json:"content"`
}

// PostMeta store meta info for blog post.
type PostMeta struct {
	Key          *datastore.Key `json:"-" datastore:"__key__"`
	Title        string         `json:"title"`
	URLSafeTitle string         `json:"url_safe_title"`
	CreatedAt    time.Time      `json:"created_at"`
	Status       PostStatus     `json:"status,string"`
	Tags         []Tag          `json:"tags"`
}

// PostContent represents blog post content data.
// All fields are indexed by default. Strings or byte slices longer than 1500 bytes cannot be indexed.
// fields used to store long strings and byte slices must be tagged with "noindex" or they will cause Put operations to fail.
type PostContent struct {
	HTML     string `json:"html" datastore:",noindex"`
	Markdown string `json:"markdown" datastore:",noindex"`
}

// PostStatus -
type PostStatus int

const (
	// PostStatusUndefined 未定義
	PostStatusUndefined PostStatus = iota
	// PostStatusPublished 公開状態
	PostStatusPublished
	// PostStatusDraft 下書き/未公開
	PostStatusDraft
	// PostStatusDeleted soft delete
	PostStatusDeleted
)

// String -
func (ps PostStatus) String() (s string) {
	switch ps {
	case PostStatusPublished:
		s = "published"
	case PostStatusDraft:
		s = "draft"
	case PostStatusDeleted:
		s = "deleted"
	default:
		s = "undefined"
	}
	return s
}

// ToPostStatus -
func ToPostStatus(s string) (ps PostStatus) {
	switch s := strings.ToLower(s); s {
	case "published":
		ps = PostStatusPublished
	case "draft":
		ps = PostStatusDraft
	case "deleted":
		ps = PostStatusDeleted
	default:
		ps = PostStatusUndefined
	}
	return ps
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

// ListMetaInput -
type ListMetaInput struct {
	Limit  int
	Offset int
	Order  PostOrder
	Status *PostStatus
}

// GetPostInput -
type GetPostInput struct {
	MetaID int64
}

// PostService is responsible for post data CRUD.
type PostService interface {
	Create(context.Context, *Post) (*Post, error)
	Update(context.Context, *Post) (*Post, error)
	ListMeta(context.Context, *ListMetaInput) ([]*PostMeta, error)
	Get(context.Context, *GetPostInput) (*Post, error)
}

// MarshalJSON -
func (ps PostStatus) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ps.String() + `"`), nil
}

// UnmarshalJSON -
func (ps *PostStatus) UnmarshalJSON(data []byte) error {
	*ps = ToPostStatus(string(data))
	return nil
}

// MarshalJSON -
// http://choly.ca/post/go-json-marshalling/
func (p *PostMeta) MarshalJSON() ([]byte, error) {
	type Alias PostMeta
	return json.Marshal(&struct {
		ID int64 `json:"id"`
		*Alias
	}{
		ID:    p.Key.ID,
		Alias: (*Alias)(p),
	})
}

// UnmarshalJSON -
func (p *PostMeta) UnmarshalJSON(data []byte) error {
	type Alias PostMeta
	aux := &struct {
		ID int64 `json:"id"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if p.Key == nil {
		p.Key = &datastore.Key{}
	}
	p.Key.ID = aux.ID
	return nil
}
