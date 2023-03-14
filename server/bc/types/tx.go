package types

import (
	"encoding/json"
	"time"
)

type Tx struct {
	Height    int64           `json:"height"`
	TxHash    string          `json:"txhash"`
	Tx        json.RawMessage `json:"tx"`
	TimeStamp time.Time       `json:"timestamp"`
	Code      int64           `json:"code"`
	Sender    string          `json:"sender"`
}

type TxResponse struct {
	Txs   []Tx  `json:"txs"`
	Total int64 `json:"total"`
}
