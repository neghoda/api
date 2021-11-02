package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/erp/api/src/cron"
	"github.com/erp/api/src/repo/ssga"
	"github.com/erp/api/src/server/handlers"
	"github.com/erp/api/src/service"

	log "github.com/sirupsen/logrus"

	nethttp "net/http"

	"github.com/erp/api/src/config"
	"github.com/erp/api/src/server/http"
	"github.com/erp/api/src/storage/postgres"
)

func main() {
	// read service cfg from os env
	cfg, err := config.Read()
	if err != nil {
		panic(err)
	}

	// init logger
	initLogger(cfg.LogLevel)

	log.Info("Service starting...")

	// prepare main context
	ctx, cancel := context.WithCancel(context.Background())
	setupGracefulShutdown(cancel)

	persistenceDB := postgres.New()
	err = persistenceDB.NewConn(&cfg.PostgresCfg)
	if err != nil {
		log.Fatal(fmt.Sprintf("cannot connect to the Postgres server %v", err))
	}

	var (
		wg         = &sync.WaitGroup{}
		ssgaClient = ssga.NewClient(
			nethttp.Client{},
		)
	)

	srv := service.New(
		&cfg,
		persistenceDB,
		ssgaClient,
	)

	crowWrapper := cron.NewCronWrapper(
		&cfg,
		persistenceDB,
		ssgaClient,
	)

	err = crowWrapper.Setup()
	if err != nil {
		log.Fatal(fmt.Sprintf("cannot setup cron jobs: %v", err))
	}

	httpSrv, err := http.New(
		&cfg.HTTPConfig,
		handlers.NewAuthHandler(srv),
		handlers.NewFundHandler(srv),
	)

	if err != nil {
		log.WithError(err).Fatal("http server init")
	}

	httpSrv.Run(ctx, wg)

	// wait while services work
	wg.Wait()
	log.Info("Service stopped")
}

func initLogger(logLevel string) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stderr)

	switch strings.ToLower(logLevel) {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "trace":
		log.SetLevel(log.TraceLevel)
	default:
		log.SetLevel(log.DebugLevel)
	}
}

func setupGracefulShutdown(stop func()) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChannel
		log.Error("Got Interrupt signal")
		stop()
	}()
}
