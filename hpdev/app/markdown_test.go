package app_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ymgyt/happy-developing/hpdev/app"
)

func TestMarkdown_Parse(t *testing.T) {
	tests := []struct {
		desc string
		md   []byte
		want []byte
	}{
		{
			desc: "just title and content",
			md: []byte(`# Title

hello gopher
`),
			want: []byte(`<h1>Title</h1>

<p>hello gopher</p>
`),
		},
		{
			desc: "code literal",
			md: []byte("```go" + `
func Hello() string {
	return "hello"
}
` + "```"),
			want: []byte(""),
		},
	}

	md := &app.Markdown{}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := md.ConvertHTML(tc.md)
			// 見やすいようにstringで比較する.
			if diff := cmp.Diff(string(got), string(tc.want)); diff != "" {
				t.Errorf("(-got +want)\n%s", diff)
			}
		})
	}
}

func TestSyntaxhighlight(t *testing.T) {
	tests := []struct {
		desc string
		htm  []byte
		want []byte
	}{
		{
			desc: "go",
			// htmlから必要..?
			htm: []byte(`<pre><code class="language-go">func Hello() string {
	return &quot;hello&quot;
}
</code></pre>
`),
			want: []byte(""),
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := app.SyntaxHighlight(tc.htm)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(string(got), string(tc.want)); diff != "" {
				t.Errorf("(-got +want)\n%s", diff)
			}
		})
	}
}
