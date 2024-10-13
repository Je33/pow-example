package static

import (
	"pow-example/internal/pkg/common"
	"pow-example/pkg/logger"
)

type Repo struct {
	db  map[int]*common.Quote
	log logger.Logger
}

func New(log logger.Logger) *Repo {
	db := make(map[int]*common.Quote)

	return &Repo{
		db:  db,
		log: log,
	}
}
