package netecho

import (
	"fmt"
	"io"
	"pow-example/pkg/errs"
	"pow-example/pkg/logger"
	"pow-example/pkg/netconn"
)

type ClientConfig struct {
	Network string
	Address string
}

type Client struct {
	config ClientConfig
	conn   netconn.Connector
	log    logger.Logger
}

func NewClient(config ClientConfig, log logger.Logger) *Client {
	return &Client{
		config: config,
		conn:   netconn.NewWithAddr(config.Network, config.Address),
		log:    log,
	}
}

func (c *Client) Connect() error {
	err := c.conn.Connect()
	if err != nil {
		return errs.New(fmt.Errorf("failed to connect: %w", err)).Log(c.log)
	}
	return nil
}

func (c *Client) Command(cmd string, req []byte) ([]byte, error) {
	err := c.conn.Write([]byte(cmd + sep + string(req)))
	if err != nil {
		return nil, errs.New(fmt.Errorf("failed to write to connection: %w", err)).Log(c.log)
	}

	respBytes, err := c.conn.Read()
	if err != nil {
		if errs.Is(err, io.EOF) {
			return nil, err
		}
		return nil, errs.New(fmt.Errorf("failed to read from connection: %w", err)).Log(c.log)
	}

	return respBytes, nil
}

func (c *Client) Close() error {
	err := c.conn.Close()
	if err != nil {
		return errs.New(fmt.Errorf("failed to close connection: %w", err)).Log(c.log)
	}
	return nil
}
