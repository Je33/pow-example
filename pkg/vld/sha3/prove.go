package sha3

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
	"strings"
)

func (p *Pow) Prove(challenge string) (string, error) {
	var nonce int
	for {
		nonceStr := fmt.Sprintf("%d", nonce)
		data := challenge + nonceStr

		// Compute SHA-3 hash (256 bits)
		hash := sha3.New256()
		_, err := hash.Write([]byte(data))
		if err != nil {
			return "", err
		}
		hashBytes := hash.Sum(nil)

		// Convert the hash to hexadecimal
		hashStr := hex.EncodeToString(hashBytes)

		// Check if the hash has the required number of leading zeros
		if strings.HasPrefix(hashStr, strings.Repeat("0", int(p.d))) {
			return nonceStr, nil
		}
		nonce++
	}
}
