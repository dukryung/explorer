package util

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func ConsPubKeyToBech32(pub []byte) string {
	return sdk.ConsAddress(ed25519.PubKey(pub).Address().Bytes()).String()
}

func ConsAddressToBech32(address []byte) string {
	return sdk.ConsAddress(address).String()
}
