package types

import (
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

type Block struct {
	Moniker          string        `json:"moniker"`
	ValidatorAddress string        `json:"val_address"`
	ConsensusAddress string        `json:"cons_address"`
	DiffTime         int64         `json:"diff_time"`
	TmBlock          tmproto.Block `json:"tm_block"`
}

type BlockResponse struct {
	Blocks []Block `json:"blocks"`
	Total  int64   `json:"total"`
}
