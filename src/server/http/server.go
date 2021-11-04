package http

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/neghoda/api/src/config"
	"github.com/neghoda/api/src/server/handlers"
	middleware "github.com/neghoda/api/src/server/http/middlewares"
	"github.com/neghoda/api/src/service"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	swagger "github.com/swaggo/http-swagger"

	healthcheck "github.com/neghoda/api/src/server/health-check"
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
	auh  *handlers.AuthHandler
	fund *handlers.FundHandler
}

func New(cfg *config.HTTP, srv *service.Service) (*Server, error) {
	httpSrv := http.Server{
		Addr: fmt.Sprintf(":%d", cfg.Port),
	}

	// build Server
	server := Server{
		config: cfg,
		auh:    handlers.NewAuthHandler(srv),
		fund:   handlers.NewFundHandler(srv),
	}

	server.setupHTTP(&httpSrv)

	return &server, nil
}

func (s *Server) setupHTTP(srv *http.Server) {
	srv.Handler = s.buildHandler()
	s.http = srv
}

// nolint: funlen,lll
func (s *Server) buildHandler() http.Handler {
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
	v1Router.Handle("/tickers", privateChain.ThenFunc(s.fund.GetTickerList)).Methods(http.MethodGet)
	v1Router.Handle("/funds", privateChain.ThenFunc(s.fund.GetFundByTicker)).Methods(http.MethodGet)

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
	}).Handler(router)
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
