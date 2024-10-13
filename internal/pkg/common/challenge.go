package common

type Challenge struct {
	Hash       string `json:"text"`
	Difficulty int64  `json:"difficulty"`
}
