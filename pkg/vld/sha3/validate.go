package sha3

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
	"strings"
)

func (p *Pow) Validate(challenge, nonce string) error {
	data := challenge + nonce

	// Compute SHA-3 hash (256 bits)
	hash := sha3.New256()
	hash.Write([]byte(data))
	hashBytes := hash.Sum(nil)

	// Convert the hash to hexadecimal
	hashStr := hex.EncodeToString(hashBytes)

	// Check if the hash has the required number of leading zeros
	if !strings.HasPrefix(hashStr, strings.Repeat("0", int(p.d))) {
		return fmt.Errorf("invalid hash: %s", hashStr)
	}

	return nil
}
