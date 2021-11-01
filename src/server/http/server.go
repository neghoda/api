package http

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/erp/api/src/config"
	"github.com/erp/api/src/server/handlers"
	middleware "github.com/erp/api/src/server/http/middlewares"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	swagger "github.com/swaggo/http-swagger"

	healthcheck "github.com/erp/api/src/server/health-check"
)

const (
	version1 = "/v1"
)

type Server struct {
	http      *http.Server
	runErr    error
	readiness bool

	config *config.HTTP

	// handlers
	auh *handlers.AuthHandler
}

func New(cfg *config.HTTP,
	authHandler *handlers.AuthHandler,
) (*Server, error) {
	httpSrv := http.Server{
		Addr: fmt.Sprintf(":%d", cfg.Port),
	}

	// build Server
	srv := Server{
		config: cfg,
		auh:    authHandler,
	}

	if err := srv.setupHTTP(&httpSrv); err != nil {
		return nil, err
	}

	return &srv, nil
}

func (s *Server) setupHTTP(srv *http.Server) error {
	handler, err := s.buildHandler()
	if err != nil {
		return err
	}

	srv.Handler = handler
	s.http = srv

	return nil
}

// nolint: funlen,lll
func (s *Server) buildHandler() (http.Handler, error) {
	var (
		router        = mux.NewRouter()
		serviceRouter = router.PathPrefix(s.config.URLPrefix).Subrouter()
		v1Router      = serviceRouter.PathPrefix(version1).Subrouter()

		publicChain  = alice.New()
		privateChain = publicChain.Append(middleware.Auth)
	)

	// public routes
	v1Router.Handle("/health", publicChain.ThenFunc(healthcheck.Health)).Methods(http.MethodGet)
	v1Router.Handle("/login", publicChain.ThenFunc(s.auh.Login)).Methods(http.MethodPost)
	v1Router.Handle("/token", publicChain.ThenFunc(s.auh.TokenRefresh)).Methods(http.MethodPost)
	v1Router.Handle("/signup", publicChain.ThenFunc(s.auh.SignUp)).Methods(http.MethodPost)

	// private routes
	v1Router.Handle("/logout", privateChain.ThenFunc(s.auh.Logout)).Methods(http.MethodDelete)

	// ================================= Swagger =================================================

	if s.config.SwaggerEnable {
		router.
			PathPrefix("/swagger/static").
			Handler(http.StripPrefix("/swagger/static", http.FileServer(http.Dir(s.config.SwaggerServeDir))))
		router.
			PathPrefix("/swagger").
			Handler(swagger.Handler(swagger.URL("/swagger/static/swagger.json")))
	}

	return cors.New(cors.Options{
		AllowedOrigins:     []string{s.config.CORSAllowedHost},
		AllowedMethods:     []string{http.MethodHead, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowedHeaders:     []string{"*"},
		AllowCredentials:   true,
		OptionsPassthrough: false,
	}).Handler(router), nil
}

func (s *Server) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	log.Info("http service: begin run")

	go func() {
		defer wg.Done()
		log.Debug("http service: addr=", s.http.Addr)
		err := s.http.ListenAndServe()
		s.runErr = err
		log.Info("http service: end run > ", err)
	}()

	go func() {
		<-ctx.Done()
		sdCtx, _ := context.WithTimeout(context.Background(), 5*time.Second) // nolint
		err := s.http.Shutdown(sdCtx)

		if err != nil {
			log.Info("http service shutdown (", err, ")")
		}
	}()

	s.readiness = true
}
