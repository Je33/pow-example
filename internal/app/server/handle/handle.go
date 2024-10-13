package handle

import (
	"context"
	"pow-example/internal/pkg/common"
	"pow-example/pkg/logger"
)

type QuoteService interface {
	GetQuote(ctx context.Context) (*common.Quote, error)
}

type PowService interface {
	Challenge() (*common.Challenge, error)
	Validate(challenge, nonce string) error
}

type Handle struct {
	quoteSrv QuoteService
	powSrv   PowService
	log      logger.Logger
}

func New(quoteSrv QuoteService, powSrv PowService, log logger.Logger) *Handle {
	return &Handle{
		quoteSrv: quoteSrv,
		powSrv:   powSrv,
		log:      log,
	}
}
