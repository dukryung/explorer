package controller

import (
	"math"

	"github.com/hessegg/nikto-explorer/server/api/event/codec"
	"github.com/hessegg/nikto-explorer/server/bc"
	bctypes "github.com/hessegg/nikto-explorer/server/bc/types"
	"github.com/hessegg/nikto-explorer/types"
)

type Tx struct{}

func init() {
	codec.RegisterCodec("/tx/latest", struct {
		Limit float64 `field:"limit" order:"0"`
		Page  float64 `field:"page" order:"1"`
	}{})
	codec.RegisterCodec("/tx/height", struct {
		Height float64 `field:"height" order:"0"`
		Limit  float64 `field:"limit" order:"1"`
		Page   float64 `field:"page" order:"2"`
	}{})
	codec.RegisterCodec("/tx/hash", struct {
		Hash string `field:"hash" order:"0"`
	}{})
	codec.RegisterCodec("/tx/address", struct {
		Address string  `field:"address" order:"0"`
		Limit   float64 `field:"limit" order:"1"`
		Page    float64 `field:"page" order:"2"`
	}{})
}

// Latest godoc
// @Summary      Get latest txs by count
// @Description  get latest txs by count
// @Tags         tx
// @Accept       json
// @Produce      json
// @Success      200  {object}  bctypes.TxResponse
// @Param        limit   query      int  true  "limit"
// @Param        page    query      int  true  "page"
// @Router       /api/tx/latest [get]
func (c Tx) Latest(ctx types.Context, limit, page float64) (bctypes.TxResponse, error) {
	paginate := bctypes.NewPaginate(limit, page)
	ctx.Logger().Debug("Tx.Latest", paginate)

	var txs []bctypes.Tx
	var total int64
	var err error

	var stats bctypes.Stats
	//total = Cache.GetCache(StatsCache).(bctypes.Stats).TxTotal
	err = Cache.GetRedisCache(ctx, StatsCache, &stats)
	if err != nil {
		total = 0
	} else {
		total = stats.TxTotal
	}

	if paginate.Offset() != 0 {
		txs, err = bc.TxProvider.GetLatestTxs(ctx.DB(), paginate)
		if err != nil {
			return bctypes.TxResponse{}, err
		}
	} else {
		err := Cache.GetRedisCache(ctx, TxCache, &txs)
		if err != nil {
			txs, err = bc.TxProvider.GetLatestTxs(ctx.DB(), paginate)
		} else {
			if paginate.Limit() < len(txs) {
				txs = txs[:paginate.Limit()]
			}
		}
	}

	return bctypes.TxResponse{
		Txs:   txs,
		Total: total,
	}, nil
}

// Height godoc
// @Summary      Get txs by height
// @Description  get txs by height
// @Tags         tx
// @Accept       json
// @Produce      json
// @Success      200  {object}  bctypes.TxResponse
// @Param        height   query      int  true  "height"
// @Param        limit    query      int  true  "limit"
// @Param        page     query      int  true  "page"
// @Router       /api/tx/hash [get]
func (c Tx) Height(ctx types.Context, height, limit, page float64) (bctypes.TxResponse, error) {
	h := int64(math.Round(height))
	paginate := bctypes.NewPaginate(limit, page)
	ctx.Logger().Debug("Tx.Height", h)

	txs, err := bc.TxProvider.GetTxsByHeight(ctx.DB(), h, paginate)
	if err != nil {
		return bctypes.TxResponse{}, err
	}

	total, err := bc.TxProvider.GetTxsCountByHeight(ctx.DB(), h)
	if err != nil {
		return bctypes.TxResponse{}, err
	}

	return bctypes.TxResponse{
		Txs:   txs,
		Total: total,
	}, nil
}

// Hash godoc
// @Summary      Get tx by hash
// @Description  get tx by hash
// @Tags         tx
// @Accept       json
// @Produce      json
// @Success      200  {object}  bctypes.Tx
// @Param        hash   query      string  true  "hash"
// @Router       /api/tx/hash [get]
func (c Tx) Hash(ctx types.Context, hash string) (bctypes.Tx, error) {
	ctx.Logger().Debug("Tx.Hash", hash)

	return bc.TxProvider.GetTxByHash(ctx.DB(), hash)
}

// Address godoc
// @Summary      Get txs by address
// @Description  get txs by address
// @Tags         tx
// @Accept       json
// @Produce      json
// @Success      200  {object}  bctypes.TxResponse
// @Param        address   query      string  true  "address"
// @Param        limit     query      int     true  "limit"
// @Param        page      query      int     true  "page"
// @Router       /api/tx/address [get]
func (c Tx) Address(ctx types.Context, address string, limit, page float64) (bctypes.TxResponse, error) {
	paginate := bctypes.NewPaginate(limit, page)
	ctx.Logger().Debug("Tx.Address", address)

	txs, err := bc.TxProvider.GetTxsByAddress(ctx.DB(), address, paginate)
	if err != nil {
		return bctypes.TxResponse{}, err
	}

	total, err := bc.TxProvider.GetTxsCountByAddress(ctx.DB(), address)
	if err != nil {
		return bctypes.TxResponse{}, err
	}

	return bctypes.TxResponse{
		Txs:   txs,
		Total: total,
	}, nil
}
