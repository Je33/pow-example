package quote

import (
	"context"
	"pow-example/internal/pkg/common"
)

func (s *Service) GetQuote(ctx context.Context) (*common.Quote, error) {
	return s.db.GetRandQuote(ctx)
}
