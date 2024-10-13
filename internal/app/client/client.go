package client

import (
	"encoding/json"
	"fmt"
	"pow-example/internal/app/client/config"
	"pow-example/internal/pkg/common"
	"pow-example/pkg/errs"
	"pow-example/pkg/logger"
	"pow-example/pkg/netecho"
	"pow-example/pkg/vld"
)

type Client struct {
	config    config.Config
	client    *netecho.Client
	validator vld.Validator
	log       logger.Logger
}

func New(conf config.Config, validator vld.Validator, log logger.Logger) *Client {
	return &Client{
		config: conf,
		client: netecho.NewClient(netecho.ClientConfig{
			Network: conf.ServerNetwork,
			Address: conf.ServerAddress,
		}, log),
		validator: validator,
		log:       log,
	}
}

func (c *Client) GetQuoteSequence() error {
	err := c.Connect()
	if err != nil {
		return errs.New(fmt.Errorf("client connect error: %w", err)).Log(c.log)
	}

	defer func() {
		err = c.Close()
		if err != nil {
			c.log.Error(fmt.Sprintf("client close error: %v", err))
		}
	}()

	err = c.GetQuote()
	if err != nil {
		return errs.New(fmt.Errorf("client start error: %w", err)).Log(c.log)
	}

	return nil
}

func (c *Client) Connect() error {
	// connect to the server
	err := c.client.Connect()
	if err != nil {
		return errs.New(fmt.Errorf("failed to connect to server: %w", err)).Log(c.log)
	}

	return nil
}

func (c *Client) Close() error {
	err := c.client.Close()
	if err != nil {
		return errs.New(fmt.Errorf("failed to close connection: %w", err)).Log(c.log)
	}
	return nil
}

func (c *Client) GetQuote() error {
	// get challenge
	mess, err := c.client.Command("challenge", nil)
	if err != nil {
		return errs.New(fmt.Errorf("failed to read message from provider: %w", err)).Log(c.log)
	}

	c.log.Info(fmt.Sprintf("message from provider: %s", mess))

	var challenge common.Challenge

	// unmarshal challenge
	err = json.Unmarshal(mess, &challenge)
	if err != nil {
		return errs.New(fmt.Errorf("failed to unmarshal challenge: %w", err)).Log(c.log)
	}

	c.log.Info(fmt.Sprintf("challenge: %s, difficulty: %d", challenge.Hash, challenge.Difficulty))

	// change difficulty should match between server and client
	if challenge.Difficulty != c.validator.Difficulty() {
		return errs.New(fmt.Errorf("difficulty does not match server: %d, client: %d", challenge.Difficulty, c.config.Difficulty)).Log(c.log)
	}

	// calculate nonce by challenge
	nonce, err := c.validator.Prove(challenge.Hash)
	if err != nil {
		return errs.New(fmt.Errorf("failed to prove challenge: %w", err)).Log(c.log)
	}

	c.log.Info(fmt.Sprintf("proved challenge: %s, nonce: %s", challenge.Hash, nonce))

	proof := common.Proof{
		Nonce:      nonce,
		Hash:       challenge.Hash,
		Difficulty: challenge.Difficulty,
	}
	proveBytes, err := json.Marshal(proof)
	if err != nil {
		return errs.New(fmt.Errorf("failed to marshal proof: %w", err)).Log(c.log)
	}

	// send nonce back
	quoteBytes, err := c.client.Command("quote", proveBytes)
	if err != nil {
		return errs.New(fmt.Errorf("failed to send proof: %w", err)).Log(c.log)
	}

	c.log.Info(fmt.Sprintf("quote received: %s", quoteBytes))

	var quote common.Quote

	// unmarshal quote
	err = json.Unmarshal(quoteBytes, &quote)
	if err != nil {
		return errs.New(fmt.Errorf("failed to unmarshal quote: %w", err)).Log(c.log)
	}

	// print quote to logs
	c.log.Info(fmt.Sprintf("quote decoded: %s, author: %s", quote.Text, quote.Author))

	return nil
}
