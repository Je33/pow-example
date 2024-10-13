package server

import (
	"context"
	"pow-example/internal/app/server/config"
	"pow-example/pkg/logger"
	"pow-example/pkg/netecho"
)

type Server struct {
	config config.Config
	server *netecho.Server
	log    logger.Logger
}

func New(conf config.Config, log logger.Logger) *Server {
	return &Server{
		config: conf,
		server: netecho.NewServer(netecho.ServerConfig{
			Network: conf.ServerNetwork,
			Address: conf.ServerAddress,
			MaxConn: conf.MaxClients,
		}, log),
		log: log,
	}
}

func (s *Server) Handle(prefix string, fun func(ctx context.Context, req []byte) ([]byte, error)) {
	s.server.Handle(prefix, fun)
}

func (s *Server) Start(ctx context.Context) {
	s.server.Listen(ctx)
}
