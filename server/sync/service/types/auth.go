package types

import (
	"encoding/base64"
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type TxAuthInfo struct {
	AuthInfo struct {
		SignerInfo []struct {
			PublicKey struct {
				Key string `json:"key"`
			} `json:"public_key"`
		} `json:"signer_infos"`
	} `json:"auth_info"`
}

func ParseSignerAddress(rawTx json.RawMessage) (sdk.AccAddress, error) {
	authInfo := TxAuthInfo{}

	err := json.Unmarshal(rawTx, &authInfo)
	if err != nil {
		return nil, err
	}

	//TODO just single signer
	decodedPubKey, err := base64.StdEncoding.DecodeString(authInfo.AuthInfo.SignerInfo[0].PublicKey.Key)
	if err != nil {
		return nil, err
	}

	pubKey := secp256k1.PubKey(decodedPubKey)
	return sdk.AccAddress(pubKey.Address().Bytes()), nil
}

