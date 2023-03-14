package types

import "time"

type Stats struct {
	TxTotal           int64  `json:"tx_total"`
	TokenTotal        int64  `json:"token_total"`
	NftTotal          int64  `json:"nft_total"`
	ValidatorTotal    int64  `json:"validator_total"`
	TotalBondedTokens string `json:"total_bonded_tokens"`
	BlockHeight       int64  `json:"block_height"`
	BlockAvgTime      int64  `json:"block_avg_time"`
	BlockMinTime      int64  `json:"block_min_time"`
}

type TxStats struct {
	TxCount   int       `json:"tx_count"`
	TimeStamp time.Time `json:"time_stamp"`
}
