package static

import (
	"context"
	"fmt"
	"math/rand"
	"pow-example/internal/pkg/common"
	"pow-example/pkg/errs"
)

func (r *Repo) GetRandQuote(ctx context.Context) (*common.Quote, error) {
	count := len(r.db)
	if count == 0 {
		return nil, errs.New(fmt.Errorf("no quotes in db")).Log(r.log)
	}

	randId := rand.Intn(count)
	return r.db[randId], nil
}
