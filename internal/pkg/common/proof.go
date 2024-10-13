package common

type Proof struct {
	Nonce      string `json:"nonce"`
	Hash       string `json:"hash"`
	Difficulty int64  `json:"difficulty"`
}
