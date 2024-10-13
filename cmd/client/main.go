package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"pow-example/internal/app/client"
	"pow-example/internal/app/client/config"
	"pow-example/pkg/cfg"
	"pow-example/pkg/logger"
	"pow-example/pkg/vld"
	"pow-example/pkg/vld/sha3"
	"sync"
	"syscall"
)

func main() {
	conf, err := cfg.Read[config.Config]()
	if err != nil {
		fmt.Println(fmt.Errorf("failed to read config: %w", err))
		return
	}

	log := logger.New(conf.LogLevel)

	ctx, cancel := context.WithCancel(context.Background())

	log.Info(fmt.Sprintf("client initialized with config: %+v", conf))

	log.Info("initializing dependencies")

	log.Info("starting client")

	algorithm := sha3.New(conf.Difficulty)
	powValidator := vld.New(algorithm)

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)

	wg := sync.WaitGroup{}

	go func() {
		for {
			getQuotesTimes := 12

			log.Info(fmt.Sprintf("starting new connection and getting quotes %d times", getQuotesTimes))

			wg = sync.WaitGroup{}
			wg.Add(getQuotesTimes)

			for i := 0; i < getQuotesTimes; i++ {
				go func() {
					defer wg.Done()

					quoteClient := client.New(conf, powValidator, log)

					err := quoteClient.GetQuoteSequence()
					if err != nil {
						log.Error(fmt.Sprintf("failed to get quote: %v", err))
					}
				}()
			}
			wg.Wait()
			if ctx.Err() != nil {
				log.Info("all processes are done, shutting down")
				return
			}
		}
	}()

	sig := <-cancelChan
	log.Info(fmt.Sprintf("shutdown signal received (%s), waiting for all processes will be done to close clients", sig.String()))

	cancel()
	wg.Wait()
}
