package cfg

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

func Read[T any]() (T, error) {
	var cfg T

	err := envconfig.Process("myapp", &cfg)
	if err != nil {
		return cfg, fmt.Errorf("failed to init config: %w", err)
	}

	return cfg, nil
}
