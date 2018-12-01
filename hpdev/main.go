package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/ymgyt/happy-developing/hpdev/gcp"

	"cloud.google.com/go/datastore"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"

	"github.com/ymgyt/happy-developing/hpdev/app"
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

	mws := middlewares.NewChain(r, &middlewares.Logging{})

	r.GET("/static/*filepath", hs.Static.ServeStatic)
	r.GET("/example", hs.Example.RenderExample)

	r.GET("/author/posts/new", hs.Post.RenderPostForm)
	r.POST("/author/posts/new", hs.Post.CreatePost)
	r.GET("/author/posts", hs.Post.ListPosts)

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

	return &app.Services{
		PostService: postService,
	}
}

func newEnv() *app.Env {
	// create logger
	log := newLogger()

	return &app.Env{
		Log: log,
		Ctx: context.Background(),
		Now: app.Now,
	}
}

func newLogger() *logrus.Logger {
	var formatter logrus.Formatter
	var level logrus.Level
	switch appMode {
	case app.DevelopmentMode:
		formatter = &logrus.TextFormatter{}
		level = logrus.DebugLevel
	default:
		formatter = &logrus.JSONFormatter{}
		level = logrus.InfoLevel
	}

	log := logrus.New()
	log.Formatter = formatter
	log.Level = level

	return log
}

func main() {
	checkEnvironments()

	env := newEnv()
	mux := registerHandlers(env, httprouter.New())

	s := server.Must(server.Config{
		Addr:    ":" + port,
		Handler: mux,
	})

	env.Log.Info("running on ", port)
	env.Log.Info(s.Run())
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
