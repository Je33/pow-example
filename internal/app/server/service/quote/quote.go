package quote

import (
	"context"
	"pow-example/internal/pkg/common"
	"pow-example/pkg/logger"
)

type DbRepo interface {
	GetRandQuote(ctx context.Context) (*common.Quote, error)
}

type Service struct {
	db  DbRepo
	log logger.Logger
}

func New(db DbRepo, log logger.Logger) *Service {
	return &Service{
		db:  db,
		log: log,
	}
}
