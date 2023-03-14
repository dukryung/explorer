package controller

import (
	"fmt"

	"github.com/hessegg/nikto-explorer/server/api/event/codec"
	"github.com/hessegg/nikto-explorer/server/bc"
	bctypes "github.com/hessegg/nikto-explorer/server/bc/types"
	"github.com/hessegg/nikto-explorer/types"
)

type Nft struct{}

func init() {
	codec.RegisterCodec("/nft/list", struct {
		Limit float64 `field:"limit" order:"0"`
		Page  float64 `field:"page" order:"1"`
	}{})
	codec.RegisterCodec("/nft/info", struct {
		Id string `field:"id" order:"0"`
	}{})

}

// List godoc
// @Summary      Get latest nft by alphabet
// @Description  get latest nft by alphabet
// @Tags         nft
// @Accept       json
// @Produce      json
// @Param        limit   query      int  true  "limit"
// @Param        page    query      int  true  "page"
// @Router       /api/nft/list [get]
func (c Nft) List(ctx types.Context, limit, page float64) (bctypes.NftInfoList, error) {
	ctx.Logger().Debug("NFT.List", limit, page)
	paginate := bctypes.NewPaginate(limit, page)
	var err error
	var stats bctypes.Stats
	var total int64
	var infos []bctypes.NftInfo

	err = Cache.GetRedisCache(ctx, StatsCache, &stats)
	if err != nil {
		total = 0
	} else {
		total = stats.NftTotal
	}
	// Try cache first
	if paginate.Offset() == 0 {
		err := Cache.GetRedisCache(ctx, NftCache, &infos)
		if err == nil && paginate.Limit() < len(infos) {
			infos = infos[:paginate.Limit()]
			return bctypes.NftInfoList{
				Infos: infos,
				Total: total,
			}, nil
		}
	}
	// fallback
	infos, err = c.list(ctx, limit, page)
	if err != nil {
		return bctypes.NftInfoList{}, nil
	}
	return bctypes.NftInfoList{
		Infos: infos,
		Total: total,
	}, nil
}

func (c Nft) ListCollections(ctx types.Context, limit, page float64) (bctypes.NftCollectionInfoList, error) {
	paginate := bctypes.NewPaginate(limit, page)

	var infos []bctypes.NftCollectionInfo
	var err error

	infos, err = bc.NftProvider.GetNftCollections(ctx.DB(), paginate)
	if err != nil {
		return bctypes.NftCollectionInfoList{}, nil
	}
	total := bc.NftProvider.GetNftCollectionCount(ctx.DB())
	return bctypes.NftCollectionInfoList{
		Infos: infos,
		Total: total,
	}, nil
}

// ListByCollection godoc
// @Summary      List NF tokens by collection id
// @Description  List NF tokens by collection id
// @Tags         nft
// @Accept       json
// @Produce      json
// @Param        collectionId   query      string  true  "id"
// @Router       /api/nft/listbycollection [get]
func (c Nft) ListByCollection(ctx types.Context, collectionId string, limit, page float64) (bctypes.NftInfoList, error) {
	ctx.Logger().Debug("NFT.ListByCollection", collectionId, limit, page)

	return c.listBy(ctx, "collection_id", collectionId, limit, page)
}

// ListByOwner godoc
// @Summary      List NF tokens by the owner address
// @Description  List NF tokens by the owner address
// @Tags         nft
// @Accept       json
// @Produce      json
// @Param        address   query      string  true  "id"
// @Router       /api/nft/listbyowner [get]
func (c Nft) ListByOwner(ctx types.Context, address string, limit, page float64) (bctypes.NftInfoList, error) {
	ctx.Logger().Debug("NFT.ListByOwner", address, limit, page)
	return c.listBy(ctx, "owner_address", address, limit, page)
}

// ListByIssuer godoc
// @Summary      List NF tokens by the issuer address
// @Description  List NF tokens by the issuer address
// @Tags         nft
// @Accept       json
// @Produce      json
// @Param        address   query      string  true  "id"
// @Router       /api/nft/listbyissuer [get]
func (c Nft) ListByIssuer(ctx types.Context, address string, limit, page float64) (bctypes.NftInfoList, error) {
	ctx.Logger().Debug("NFT.ListByIssuer", address, limit, page)
	return c.listBy(ctx, "issuer_address", address, limit, page)
}

// ListByDelegatedTo godoc
// @Summary      List NF tokens by the delegated-to address
// @Description  List NF tokens by the delegated-to address
// @Tags         nft
// @Accept       json
// @Produce      json
// @Param        address   query      string  true  "id"
// @Router       /api/nft/listbydelegatedto [get]
func (c Nft) ListByDelegatedTo(ctx types.Context, address string, limit, page float64) (bctypes.NftInfoList, error) {
	ctx.Logger().Debug("NFT.ListByDelegatedTo", address, limit, page)
	return c.listBy(ctx, "delegated_to_address", address, limit, page)
}

// Info godoc
// @Summary      Get NF token by id
// @Description  get NF token by id
// @Tags         token
// @Accept       json
// @Produce      json
// @Param        id   query      string  true  "id"
// @Router       /api/nft/info [get]
func (c Nft) Info(ctx types.Context, id string) (bctypes.NftInfo, error) {
	ctx.Logger().Debug("NFT.Info", id)
	return bc.NftProvider.GetNftInfo(ctx.DB(), id)
}

func (c Nft) list(ctx types.Context, limit, page float64, where ...bc.WhereClause) ([]bctypes.NftInfo, error) {
	paginate := bctypes.NewPaginate(limit, page)

	var infos []bctypes.NftInfo
	var err error

	infos, err = bc.NftProvider.GetNftList(ctx.DB(), paginate, where...)
	if err != nil {
		return []bctypes.NftInfo{}, err
	}
	return infos, nil
}

func (c Nft) count(ctx types.Context, where ...bc.WhereClause) int64 {
	result := bc.NftProvider.GetNftCount(ctx.DB(), where...)
	ctx.Logger().Debug("nft.count", "where", where, result)
	return result
}

func (c Nft) listBy(ctx types.Context, column string, value string, limit, page float64) (bctypes.NftInfoList, error) {
	list_where_clause := bc.WhereClause{
		Clause: fmt.Sprintf(`%s = $3`, column),
		Values: []any{value},
	}
	count_where_clause := bc.WhereClause{
		Clause: fmt.Sprintf(`%s = $1`, column),
		Values: []any{value},
	}
	infos, err := c.list(ctx, limit, page, list_where_clause)
	if err != nil {
		return bctypes.NftInfoList{}, nil
	}
	return bctypes.NftInfoList{
		Infos: infos,
		Total: c.count(ctx, count_where_clause),
	}, nil
}
