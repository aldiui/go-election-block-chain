package domain

type Header struct {
	Nonce    int    `json:"nonce"`
	PrevHash string `json:"prev_hash"`
	Time     int64  `json:"time"`
}
