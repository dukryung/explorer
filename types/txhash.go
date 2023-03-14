package types

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func NewTxHash(txBytes []byte) TxHash {
	return TxHash{
		bytes: txBytes,
		hash:  sha256.Sum256(txBytes),
	}
}

type TxHash struct {
	bytes []byte
	hash  [32]byte
}

func (tx TxHash) String() string {
	return strings.ToUpper(hex.EncodeToString(tx.hash[:]))
}

func (tx TxHash) Byte() [32]byte {
	return tx.hash
}