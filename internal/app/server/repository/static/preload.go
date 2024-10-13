package static

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"pow-example/pkg/errs"
)

func (r *Repo) Preload(ctx context.Context, filePath string) error {
	r.log.Info(fmt.Sprintf("opening %s", filePath))
	file, err := os.Open(filePath)

	if err != nil {
		return errs.New(fmt.Errorf("failed to open %s: %w", filePath, err)).Log(r.log)
	}

	r.log.Info(fmt.Sprintf("successfully opened %s", filePath))

	defer func() {
		r.log.Info(fmt.Sprintf("closing %s", filePath))
		err = file.Close()
		if err != nil {
			r.log.Error(fmt.Sprintf("failed to close %s: %v", filePath, err))
		}
	}()

	decoder := json.NewDecoder(file)

	r.log.Info(fmt.Sprintf("decoding %s", filePath))
	err = decoder.Decode(&r.db)
	if err != nil {
		return errs.New(fmt.Errorf("failed to decode %s: %w", filePath, err)).Log(r.log)
	}

	r.log.Info(fmt.Sprintf("successfully decoded %s", filePath))

	return nil
}
