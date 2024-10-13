package pow

import "pow-example/pkg/logger"

type Validator interface {
	Difficulty() int64
	Generate() (string, error)
	Validate(challenge, nonce string) error
	Prove(challenge string) (string, error)
}

type Service struct {
	vld Validator
	log logger.Logger
}

func New(vld Validator, log logger.Logger) *Service {
	return &Service{
		vld: vld,
		log: log,
	}
}
