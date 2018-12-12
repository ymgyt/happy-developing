package server

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/ymgyt/happy-developing/hpdev/app"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"

	"cloud.google.com/go/datastore"
)

// New -
func New(cfg Config) (*Server, error) {
	return &Server{
		Config: cfg,
		s:      httpServer(&cfg),
	}, nil
}

// Must -
func Must(cfg Config) *Server {
	s, err := New(cfg)
	if err != nil {
		panic(err)
	}
	return s
}

// Config -
type Config struct {
	Addr            string
	DisableHTTPS    bool
	Handler         http.Handler
	DatastoreClient *datastore.Client
	Env             *app.Env
}

// Server -
type Server struct {
	Config
	s *http.Server
}

// Run -
func (s *Server) Run() error {
	if s.Config.DisableHTTPS {
		return s.s.ListenAndServe()
	}
	return s.s.ListenAndServeTLS("", "")
}

func httpServer(cfg *Config) *http.Server {
	s := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      cfg.Handler,
		Addr:         cfg.Addr,
	}

	if cfg.DisableHTTPS {
		return s
	}

	cm := certManager(cfg)
	s.TLSConfig = &tls.Config{GetCertificate: cm.GetCertificate}

	go handleAuthCallback(cfg, cm)

	return s
}

func certManager(cfg *Config) *autocert.Manager {
	m := &autocert.Manager{
		Cache:      newDatastoreCertCache(cfg.DatastoreClient),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: hostPolicy(cfg.Env),
	}
	return m
}

// TODO more strict
func hostPolicy(env *app.Env) autocert.HostPolicy {
	return func(ctx context.Context, host string) error {
		env.Log.Debug("host_policy", zap.String("host", host))
		return nil
	}
}

// Let's encryptのcallbackのhandling. 仕様をよくわかっていない
// see https://goenning.net/2017/11/08/free-and-automated-ssl-certificates-with-go/
func handleAuthCallback(cfg *Config, cm *autocert.Manager) {
	h := cm.HTTPHandler(nil)
	err := http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		dump, err := httputil.DumpRequest(r, false)
		if err == nil {
			cfg.Env.Log.Debug("autocert", zap.String("req_dump", string(dump)))
		}
	}))
	if err != nil {
		cfg.Env.Log.Error("autocert", zap.Error(err))
	}
}
