package pow

import (
	"fmt"
	"pow-example/internal/pkg/common"
	"pow-example/pkg/errs"
)

func (s *Service) Challenge() (*common.Challenge, error) {
	challenge, err := s.vld.Generate()
	if err != nil {
		return nil, errs.New(fmt.Errorf("failed to generate challenge: %w", err)).Log(s.log)
	}

	return &common.Challenge{
		Hash:       challenge,
		Difficulty: s.vld.Difficulty(),
	}, nil
}
