package netecho

import (
	"context"
	"fmt"
	"io"
	"net"
	"pow-example/pkg/errs"
	"pow-example/pkg/logger"
	"pow-example/pkg/netconn"
	"strings"
	"sync/atomic"
)

const (
	// separator for command and request body
	sep = ":::"
)

type ServerConfig struct {
	Network string
	Address string
	MaxConn int64
}

type Server struct {
	config   ServerConfig
	stats    stats
	log      logger.Logger
	handlers map[string]func(ctx context.Context, req []byte) ([]byte, error)
}

type stats struct {
	clients atomic.Int64
}

func NewServer(conf ServerConfig, log logger.Logger) *Server {
	return &Server{
		config:   conf,
		handlers: make(map[string]func(ctx context.Context, req []byte) ([]byte, error)),
		log:      log,
	}
}

func (s *Server) Handle(prefix string, fun func(ctx context.Context, req []byte) ([]byte, error)) {
	s.handlers[prefix] = fun
}

func (s *Server) process(ctx context.Context, conn netconn.Connector) ([]byte, error) {
	reqBytes, err := conn.Read()
	if err != nil {
		if errs.Is(err, io.EOF) {
			return nil, err
		}
		return nil, errs.New(fmt.Errorf("failed to read from connection: %w", err)).Log(s.log)
	}

	reqStr := string(reqBytes)
	parts := strings.SplitN(reqStr, sep, 2)

	handle, found := s.handlers[parts[0]]
	if !found {
		return nil, errs.New(fmt.Errorf("unknown command: %s", parts[0])).Log(s.log)
	}

	return handle(ctx, []byte(parts[1]))
}

func (s *Server) listenCommand(ctx context.Context, conn netconn.Connector) {
	defer func() {
		err := conn.Close()
		if err != nil {
			s.log.Error(fmt.Sprintf("failed to close connection: %v", err))
		}
		s.stats.clients.Add(-1)
	}()

	for {
		bytes, err := s.process(ctx, conn)
		if err != nil {
			if errs.Is(err, io.EOF) {
				return
			}
			s.log.Error(fmt.Sprintf("failed to process command: %v", err))
			continue
		}

		err = conn.Write(bytes)
		if err != nil {
			s.log.Error(fmt.Sprintf("failed to write to connection: %v", err))
		}
	}
}

func (s *Server) Listen(ctx context.Context) {
	listener, err := net.Listen(s.config.Network, s.config.Address)
	if err != nil {
		s.log.Error(fmt.Sprintf("failed to start server: %v", err))
		return
	}

	go func() {
		<-ctx.Done()

		s.log.Info("shutting down listener")
		err := listener.Close()
		if err != nil {
			s.log.Error(fmt.Sprintf("failed to close listener: %v", err))
		}

		s.log.Info("listener closed")
	}()

	for {
		s.log.Info("waiting new connection")
		conn, err := listener.Accept()
		if err != nil {
			if errs.Is(err, net.ErrClosed) {
				s.log.Info("server stopped")
				return
			}
			s.log.Error(fmt.Sprintf("failed to accept connection: %v", err))
			continue
		}

		curClients := s.stats.clients.Load()
		if curClients >= s.config.MaxConn {
			s.log.Error(fmt.Sprintf("max clients (%d) reached: %s", curClients, conn.RemoteAddr().String()))
			err := conn.Close()
			if err != nil {
				s.log.Error(fmt.Sprintf("failed to close connection %s: %v", conn.RemoteAddr().String(), err))
			}
			continue
		}

		cl := netconn.NewWithConn(conn)

		s.stats.clients.Add(1)

		go s.listenCommand(ctx, cl)

		s.log.Info(fmt.Sprintf("client connected: %s, total clients: %d", conn.RemoteAddr().String(), curClients+1))
	}
}
