package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ymgyt/blogo/blogo/app"

	"github.com/ymgyt/blogo/blogo/server"

	"github.com/ymgyt/blogo/blogo/handlers"

	"github.com/julienschmidt/httprouter"
)

const (
	defaultPort   = "8123"
	defaultAppEnv = "development"
)

var (
	appRoot string
	port    string
	appEnv  app.Mode
)

func registerHandlers(r *httprouter.Router) {
	hs, err := handlers.New(handlers.Config{
		AppRoot:              appRoot,
		StaticPath:           "/static",
		TemplatePath:         "/templates",
		AlwaysParseTemplates: appEnv == app.DevelopmentMode,
	})
	if err != nil {
		fail(err.Error())
	}

	r.GET("/static/*filepath", hs.Static.ServeStatic)
	r.GET("/example", hs.Example.RenderExample)

}

func checkEnvironments() {
	if appRoot == "" {
		fail("environment variable APP_ROOT required")
	}
	if appEnv == app.UndefinedMode {
		fail("environment variable APP_MODE required")
	}
}

func fail(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func main() {
	checkEnvironments()

	r := httprouter.New()
	registerHandlers(r)

	s := server.Must(server.Config{
		Addr:    ":" + port,
		Handler: r,
	})

	log.Printf("running on %s\n", port)
	log.Println(s.Run())
}

func init() {
	appRoot = os.Getenv("APP_ROOT")
	port = os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	switch m := strings.ToLower(os.Getenv("APP_MODE")); m {
	case "dev", "development":
		appEnv = app.DevelopmentMode
	case "test", "testing":
		appEnv = app.TestingMode
	case "prod", "production":
		appEnv = app.ProductionMode
	default:
		appEnv = app.DevelopmentMode
	}
}
