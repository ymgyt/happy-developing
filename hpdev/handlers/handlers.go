package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"cloud.google.com/go/datastore"

	"github.com/ymgyt/happy-developing/hpdev/app"
)

const (
	templateExt = ".html"
)

// New -
func New(cfg Config) (*Handlers, error) {
	ts, err := makeTemplateSet(filepath.Join(cfg.AppRoot, cfg.TemplatePath), cfg.AlwaysParseTemplates)
	if err != nil {
		return nil, err
	}
	base := &base{Env: cfg.Env}

	static := (&Static{base: base}).StaticRoot(path.Join(cfg.AppRoot, cfg.StaticPath), cfg.StaticPath)

	return &Handlers{
		Example:  &Example{base: base, ts: ts, templateName: "layouts/base"},
		Static:   static,
		Post:     &Post{base: base, ts: ts, templateName: "new_post", service: cfg.Services.PostService},
		Markdown: &Markdown{base: base},
		Tag:      &Tag{base: base, service: cfg.Services.TagService},
	}, nil
}

// Config -
type Config struct {
	Env                  *app.Env
	Services             *app.Services
	AppRoot              string
	StaticPath           string
	TemplatePath         string
	AlwaysParseTemplates bool
}

// Handlers -
type Handlers struct {
	Example  *Example
	Static   *Static
	Post     *Post
	Markdown *Markdown
	Tag      *Tag
}

type templateSet struct {
	reload bool
	root   string
	*template.Template
}

// ExecuteTemplate wrap inner template.Template.ExecuteTemplate() for provider reload utility.
func (ts *templateSet) ExecuteTemplate(w io.Writer, name string, data interface{}) error {
	if ts.reload {
		return ts.executeLatestTemplate(w, name, data)
	}
	return ts.Template.ExecuteTemplate(w, name, data)
}

func (ts *templateSet) executeLatestTemplate(w io.Writer, name string, data interface{}) error {
	// TODO need lock
	refreshed, err := findTemplateFiles(ts.root)
	if err != nil {
		return err
	}

	t, err := newTemplate(refreshed)
	if err != nil {
		return fmt.Errorf("failed to parse template files: %s", err)
	}

	return t.ExecuteTemplate(w, name, data)
}

//cSpell:words tmpls
func makeTemplateSet(root string, alwaysParse bool) (*templateSet, error) {
	tmpls, err := findTemplateFiles(root)
	if err != nil {
		return nil, err
	}
	t, err := newTemplate(tmpls)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template files: %s", err)
	}

	return &templateSet{Template: t, root: root, reload: alwaysParse}, nil
}

func findTemplateFiles(root string) ([]string, error) {
	var tmpls []string
	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if isTemplateFile(path) {
			tmpls = append(tmpls, path)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to walk template dirs: %s", err)
	}
	return tmpls, nil
}

func isTemplateFile(path string) bool {
	return strings.HasSuffix(path, templateExt)
}

func newTemplate(files []string) (*template.Template, error) {
	return template.New("templateSet").Funcs(funcMap).ParseFiles(files...)
}

var funcMap = template.FuncMap{
	"formatTime":        tfFormatTime,
	"formatPostMetaKey": tfFormatPostMetaID,
	"formatBytes":       tfFormatBytes,
}

func tfFormatTime(t time.Time) string {
	return t.Format("2006-01-02:15:04:05")
}

func tfFormatPostMetaID(k *datastore.Key) string {
	return strconv.FormatInt(k.ID, 10)
}

func tfFormatBytes(b []byte) string {
	return string(b)
}

// base provide common behavior to handlers.
type base struct {
	*app.Env
}

func (b base) handleRenderError(err error) {
	if err == nil {
		return
	}

	b.Env.Log.Error("render", zap.Error(err))
}

func (b *base) invalidRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
}

type apiResponse struct {
	Data interface{}
	Err  error
	W    http.ResponseWriter
}

func (b *base) json(r *apiResponse) {
	code := http.StatusOK
	enc := json.NewEncoder(r.W)
	var encErr error
	if r.Err != nil {
		code = http.StatusInternalServerError
		r.W.WriteHeader(code)
		encErr = enc.Encode(r.Err.Error())
	} else {
		r.W.WriteHeader(code)
		encErr = enc.Encode(r.Data)
	}
	if encErr != nil {
		b.Env.Log.Warn("encode", zap.Error(encErr))
	}
}
