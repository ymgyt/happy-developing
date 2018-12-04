package view

import (
	"github.com/ymgyt/happy-developing/hpdev/app"
)

// AuthorPost postの編集画面を表現するview.
type AuthorPost struct {
	Data       *app.Post
	FormAction string
}
