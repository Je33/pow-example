package sha3

type Pow struct {
	d int64
}

func New(difficulty int64) *Pow {
	return &Pow{
		d: difficulty,
	}
}

func (p *Pow) Difficulty() int64 {
	return p.d
}
