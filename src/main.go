package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/erp/api/src/server/handlers"
	"github.com/erp/api/src/service"
	"github.com/erp/api/src/storage/redis"

	log "github.com/sirupsen/logrus"

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

	var wg = &sync.WaitGroup{}

	db, err := postgres.New(ctx, wg, &cfg.PostgresCfg, &cfg.PostgresReplicaCfg)
	if err != nil {
		log.WithError(err).Fatal("postgres connection error")
	}

	redisClient, err := redis.New(&cfg.RedisCfg)
	if err != nil {
		log.WithError(err).Fatal("redis connection error")
	}

	srv := service.New(
		&cfg,
		redisClient,
		db.NewAuthRepo(),
		db.NewProfileRepo(),
	)

	httpSrv, err := http.New(
		&cfg.HTTPConfig,
		handlers.NewAuthHandler(srv),
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