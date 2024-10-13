package vld

type Validator interface {
	Generate() (string, error)
	Validate(challenge, nonce string) error
	Prove(challenge string) (string, error)
	Difficulty() int64
}

type Mechanism interface {
	Difficulty() int64
	Generate() (string, error)
	Validate(challenge, nonce string) error
	Prove(challenge string) (string, error)
}

type Validate struct {
	mechanism Mechanism
}

func New(mechanism Mechanism) *Validate {
	return &Validate{
		mechanism: mechanism,
	}
}

func (v *Validate) Generate() (string, error) {
	return v.mechanism.Generate()
}

func (v *Validate) Validate(challenge, nonce string) error {
	return v.mechanism.Validate(challenge, nonce)
}

func (v *Validate) Prove(challenge string) (string, error) {
	return v.mechanism.Prove(challenge)
}

func (v *Validate) Difficulty() int64 {
	return v.mechanism.Difficulty()
}
