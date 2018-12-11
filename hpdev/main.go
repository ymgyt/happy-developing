package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/urfave/negroni"
	"github.com/ymgyt/happy-developing/hpdev/oauth2"

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
	appRoot            string
	port               string
	appMode            app.Mode
	gcpCredentialJSON  string
	gcpProjectID       string
	githubClientID     string
	githubClientSecret string
	authorEmail        string
)

func getUrlParams(router *httprouter.Router, req *http.Request) httprouter.Params {
	_, params, _ := router.Lookup(req.Method, req.URL.Path)
	return params
}

func callwithParams(router *httprouter.Router, handler func(w http.ResponseWriter, r *http.Request, ps httprouter.Params)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := getUrlParams(router, r)
		handler(w, r, params)
	}
}

func registerHandlers(env *app.Env, r *httprouter.Router) http.Handler {
	services := newServices(env)
	hs, err := handlers.New(handlers.Config{
		Env:                  env,
		Services:             services,
		AppRoot:              appRoot,
		StaticPath:           "/static",
		TemplatePath:         "/templates",
		AlwaysParseTemplates: appMode == app.DevelopmentMode,
		OAuth2Config:         newOAuth2Config(),
	})
	if err != nil {
		fail(err.Error())
	}

	jwtMW := &middlewares.JWTVerifier{JWT: services.JWTService, Env: env}
	authorizerMW := &middlewares.Authorizer{Env: env, Email: authorEmail}

	withAuthorize := func(h httprouter.Handle) http.Handler {
		n := negroni.New(jwtMW, authorizerMW)
		n.UseHandlerFunc(callwithParams(r, h))
		return n
	}

	r.Handler("GET", "/author/posts", withAuthorize(hs.Post.RenderMetaList))
	r.Handler("GET", "/author/posts/:metaid", withAuthorize(hs.Post.RenderPostForm)) // 新規投稿の場合は metaid => new

	r.GET("/static/*filepath", hs.Static.ServeStatic)
	r.GET("/example", hs.Example.RenderExample) // TODO cleanup

	// auth
	r.GET("/login", hs.OAuth2.RenderLogin)
	r.GET("/oauth/github/callback", hs.OAuth2.GithubCallback)
	r.GET("/authenticate", hs.Auth.Authenticate)

	// api
	// authorizeに組み込む
	r.GET("/api/author/posts/:metaid", hs.Post.Get)
	r.GET("/api/author/tags", hs.Tag.List)

	r.PUT("/api/author/posts/:metaid", hs.Post.Update)

	r.POST("/api/author/posts/new", hs.Post.Create)
	r.POST("/api/author/tags", hs.Tag.Create)
	r.POST("/markdown", hs.Markdown.ConvertHTML)

	// 共通で適用するmiddleware
	// panic recoverも必要
	common := negroni.New(middlewares.MustLogging(env))
	common.UseHandler(r)

	return common
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
	jwtService := &app.JWT{HMACSecret: []byte("should_more_secret_random_value")}

	return &app.Services{
		PostService: postService,
		TagService:  tagService,
		JWTService:  jwtService,
	}
}

func newEnv() *app.Env {
	return &app.Env{
		Mode:      appMode,
		Log:       newLogger(),
		Ctx:       context.Background(),
		Now:       app.Now,
		Validator: app.MustValidator(),
	}
}

func newLogger() *zap.Logger {
	return app.MustLogger(&app.LoggingConfig{Mode: appMode, Out: os.Stdout})
}

func newOAuth2Config() *oauth2.Config {
	return &oauth2.Config{
		Github: &oauth2.Entry{
			Endpoint: &oauth2.Endpoint{
				AuthorizeURL: "https://github.com/login/oauth/authorize",
				TokenURL:     "https://github.com/login/oauth/access_token",
			},
			Credential: &oauth2.Credential{
				ClientID:     githubClientID,
				ClientSecret: githubClientSecret,
			},
			CallbackURL: "http://localhost:8123/oauth/github/callback",
		},
		CSRFToken: "should_be_random",
	}
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

// NEED environment variable manager like env config.
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
	if githubClientID == "" {
		fail("environment variable GITHUB_CLIENT_ID required")
	}
	if githubClientSecret == "" {
		fail("environment variable GITHUB_CLIENT_SECRET required")
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
	githubClientID = os.Getenv("GITHUB_CLIENT_ID")
	githubClientSecret = os.Getenv("GITHUB_CLIENT_SECRET")
	authorEmail = os.Getenv("AUTHOR_EMAIL")

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
