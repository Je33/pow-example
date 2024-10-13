package sha3

import (
	"crypto/rand"
	"encoding/hex"
)

func (p *Pow) Generate() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
