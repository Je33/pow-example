package main

import (
	"context"
	"fmt"
	"os/signal"
	"pow-example/internal/app/server"
	"pow-example/internal/app/server/config"
	"pow-example/internal/app/server/handle"
	"pow-example/internal/app/server/repository/static"
	"pow-example/internal/app/server/service/pow"
	"pow-example/internal/app/server/service/quote"
	"pow-example/pkg/cfg"
	"pow-example/pkg/logger"
	"pow-example/pkg/vld"
	"pow-example/pkg/vld/sha3"
	"syscall"
)

func main() {
	// load configuration
	conf, err := cfg.Read[config.Config]()
	if err != nil {
		fmt.Println(fmt.Errorf("failed to read config: %w", err))
		return
	}

	// init logger
	log := logger.New(conf.LogLevel)

	log.Info(fmt.Sprintf("server initialized with config: %+v", conf))

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log.Info("initializing dependencies")

	// static storage for quotes
	repo := static.New(log)
	err = repo.Preload(ctx, conf.FilePath)
	if err != nil {
		log.Error(fmt.Sprintf("failed to preload quotes: %v", err))
		return
	}

	// algorithm for proof-of-work validator
	algorithm := sha3.New(conf.Difficulty)
	// init validator with algorithm
	powValidator := vld.New(algorithm)
	// init validation service
	powService := pow.New(powValidator, log)

	// init quote service
	quoteService := quote.New(repo, log)
	// handler to get quotes through proof-of-work check
	quoteHandler := handle.New(quoteService, powService, log)
	// net server with listener
	quoteServer := server.New(conf, log)

	// handle commands
	quoteServer.Handle("challenge", quoteHandler.GetChallenge)
	quoteServer.Handle("quote", quoteHandler.GetQuote)

	// start net listener
	go quoteServer.Start(ctx)

	log.Info(fmt.Sprintf("server started on: %s", conf.ServerAddress))

	<-ctx.Done()
}
