package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankz "github.com/hessegg/nikto/x/bankz/types"
)

// Token for instead
type Token struct {
	bankz.Token
	Coin *sdk.Coin `json:"coin,omitempty"`
}

type TokenResponse struct {
	Tokens []Token `json:"tokens"`
	Total  int64   `json:"total"`
}
