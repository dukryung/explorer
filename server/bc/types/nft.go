package types

import (
	nft "github.com/hessegg/nikto/x/nft/types"
)

// NFToken information
type NftInfo struct {
	nft.NFTInfo
}

type NftInfoList struct {
	Infos []NftInfo `json:"infos"`
	Total int64     `json:"total"`
}

type NftCollectionInfo struct {
	NftInfo
	Count int64 `json:"count"`
}

type NftCollectionInfoList struct {
	Infos []NftCollectionInfo `json:"infos"`
	Total int64               `json:"total"`
}
