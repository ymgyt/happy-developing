package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"google.golang.org/api/option"

	"github.com/ymgyt/happy-developing/hpdev/app"
	"github.com/ymgyt/happy-developing/hpdev/gcp"
	"github.com/ymgyt/happy-developing/hpdev/handlers"
	"github.com/ymgyt/happy-developing/hpdev/middlewares"
	"github.com/ymgyt/happy-developing/hpdev/server"
)

const (
	defaultPort   = "8123"
	defaultAppEnv = "development"
)

var (
	appRoot           string
	port              string
	appMode           app.Mode
	gcpCredentialJSON string
	gcpProjectID      string
)

func registerHandlers(env *app.Env, r *httprouter.Router) http.Handler {
	services := newServices(env)
	hs, err := handlers.New(handlers.Config{
		Env:                  env,
		Services:             services,
		AppRoot:              appRoot,
		StaticPath:           "/static",
		TemplatePath:         "/templates",
		AlwaysParseTemplates: appMode == app.DevelopmentMode,
	})
	if err != nil {
		fail(err.Error())
	}

	mws := middlewares.NewChain(r, &middlewares.Logging{Log: env.Log})

	r.GET("/static/*filepath", hs.Static.ServeStatic)
	r.GET("/example", hs.Example.RenderExample)

	// author
	r.GET("/author/posts", hs.Post.RenderMetaList)
	r.GET("/author/posts/:metaid", hs.Post.RenderPostForm) // 新規投稿の場合は metaid => new

	// api
	r.GET("/api/author/posts/:metaid", hs.Post.Get)
	r.GET("/api/author/tags", hs.Tag.List)

	r.PUT("/api/author/posts/:metaid", hs.Post.Update)

	r.POST("/api/author/posts", hs.Post.Create)
	r.POST("/api/author/tags", hs.Tag.Create)
	r.POST("/markdown", hs.Markdown.ConvertHTML)

	return mws
}

func fail(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func newServices(env *app.Env) *app.Services {
	datastoreClient, err := datastore.NewClient(env.Ctx, gcpProjectID, option.WithCredentialsFile(gcpCredentialJSON))
	if err != nil {
		fail(err.Error())
	}
	postService, err := gcp.NewPostStore(env, datastoreClient)
	if err != nil {
		fail(err.Error())
	}
	tagService, err := gcp.NewTagStore(env, datastoreClient)
	if err != nil {
		fail(err.Error())
	}

	return &app.Services{
		PostService: postService,
		TagService:  tagService,
	}
}

func newEnv() *app.Env {
	return &app.Env{
		Log: newLogger(),
		Ctx: context.Background(),
		Now: app.Now,
	}
}

func newLogger() *zap.Logger {
	return app.MustLogger(&app.LoggingConfig{Mode: appMode, Out: os.Stdout})
}

func main() {
	checkEnvironments()

	env := newEnv()
	mux := registerHandlers(env, httprouter.New())

	s := server.Must(server.Config{
		Addr:    ":" + port,
		Handler: mux,
	})

	env.Log.Info(fmt.Sprintf("running %s mode on %s", appMode, port))
	env.Log.Error("server", zap.Error(s.Run()))
}

func checkEnvironments() {
	if appRoot == "" {
		fail("environment variable APP_ROOT required")
	}
	if appMode == app.UndefinedMode {
		fail("environment variable APP_MODE required")
	}
	if gcpProjectID == "" {
		fail("environment variable GCP_PROJECT_ID required")
	}
	if gcpCredentialJSON == "" {
		fail("environment variable GCP_CREDENTIAL_JSON required")
	}
}

func init() {
	appRoot = os.Getenv("APP_ROOT")
	port = os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	gcpProjectID = os.Getenv("GCP_PROJECT_ID")
	gcpCredentialJSON = os.Getenv("GCP_CREDENTIAL_JSON")

	switch m := strings.ToLower(os.Getenv("APP_MODE")); m {
	case "dev", "development":
		appMode = app.DevelopmentMode
	case "test", "testing":
		appMode = app.TestingMode
	case "prod", "production":
		appMode = app.ProductionMode
	default:
		appMode = app.DevelopmentMode
	}
}
