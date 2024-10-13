package test

import (
	"context"
	"github.com/stretchr/testify/require"
	"pow-example/internal/app/client"
	configClient "pow-example/internal/app/client/config"
	"pow-example/internal/app/server"
	configServer "pow-example/internal/app/server/config"
	"pow-example/internal/app/server/handle"
	"pow-example/internal/app/server/repository/static"
	"pow-example/internal/app/server/service/pow"
	"pow-example/internal/app/server/service/quote"
	"pow-example/pkg/logger"
	"pow-example/pkg/vld"
	"pow-example/pkg/vld/sha3"
	"testing"
)

func TestClientServer(t *testing.T) {
	log := logger.New(logger.LevelDebug)

	cConf := configClient.Config{
		ServerNetwork: "tcp",
		ServerAddress: ":8080",
		Difficulty:    5,
		LogLevel:      logger.LevelDebug,
	}

	sConf := configServer.Config{
		FilePath:      "testdata/quotes.json",
		ServerNetwork: "tcp",
		ServerAddress: "127.0.0.1:8080",
		MaxClients:    10,
		Difficulty:    5,
		LogLevel:      logger.LevelDebug,
	}

	t.Run("GetQuote_Success", func(t *testing.T) {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)

		defer cancel()

		algorithm := sha3.New(5)
		powVld := vld.New(algorithm)
		powSrv := pow.New(powVld, log)

		quoteClient := client.New(cConf, powVld, log)

		quoteRepo := static.New(log)

		err := quoteRepo.Preload(ctx, sConf.FilePath)
		require.NoError(t, err)

		quoteSrv := quote.New(quoteRepo, log)

		handler := handle.New(quoteSrv, powSrv, log)

		srv := server.New(sConf, log)

		srv.Handle("challenge", handler.GetChallenge)
		srv.Handle("quote", handler.GetQuote)

		go srv.Start(ctx)

		err = quoteClient.GetQuoteSequence()
		require.NoError(t, err)
	})

}
