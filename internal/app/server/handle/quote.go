package handle

import (
	"context"
	"encoding/json"
	"fmt"
	"pow-example/internal/pkg/common"
	"pow-example/pkg/errs"
)

func (h *Handle) GetQuote(ctx context.Context, req []byte) ([]byte, error) {
	h.log.Info(fmt.Sprintf("received get quote command: %s", req))

	nonce := common.Proof{}
	err := json.Unmarshal(req, &nonce)
	if err != nil {
		return nil, errs.New(fmt.Errorf("failed to unmarshal nonce (%s): %w", req, err)).Log(h.log)
	}

	h.log.Info(fmt.Sprintf("received nonce: %+v", nonce))

	// validate nonce
	err = h.powSrv.Validate(nonce.Hash, nonce.Nonce)
	if err != nil {
		return nil, errs.New(fmt.Errorf("failed to validate nonce: %w", err)).Log(h.log)
	}

	h.log.Info(fmt.Sprintf("validated nonce: %+v", nonce))

	// get quote from storage
	quote, err := h.quoteSrv.GetQuote(ctx)
	if err != nil {
		return nil, errs.New(fmt.Errorf("failed to get quote: %w", err)).Log(h.log)
	}

	h.log.Info(fmt.Sprintf("generated quote: %s", quote))

	// marshal quote to json
	decodedQuote, err := json.Marshal(quote)
	if err != nil {
		return nil, errs.New(fmt.Errorf("failed to decode quote: %w", err)).Log(h.log)
	}

	h.log.Info(fmt.Sprintf("sending quote: %s", decodedQuote))

	return decodedQuote, nil
}
