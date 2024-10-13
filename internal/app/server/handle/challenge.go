package handle

import (
	"context"
	"encoding/json"
	"fmt"
	"pow-example/pkg/errs"
)

func (h *Handle) GetChallenge(ctx context.Context, req []byte) ([]byte, error) {
	h.log.Info(fmt.Sprintf("received get challenge command: %s", req))

	chMsg, err := h.powSrv.Challenge()
	if err != nil {
		return nil, errs.New(fmt.Errorf("failed to generate challenge: %w", err)).Log(h.log)
	}

	return json.Marshal(chMsg)
}
