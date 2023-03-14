package controller

import (
	"math"

	"github.com/hessegg/nikto-explorer/server/api/event/codec"
	"github.com/hessegg/nikto-explorer/server/bc"
	bctypes "github.com/hessegg/nikto-explorer/server/bc/types"
	"github.com/hessegg/nikto-explorer/types"
)

type Block struct{}

func init() {
	codec.RegisterCodec("/block/latest", struct {
		Limit float64 `field:"limit" order:"0"`
		Page  float64 `field:"page" order:"1"`
	}{})
	codec.RegisterCodec("/block/height", struct {
		Height float64 `field:"height" order:"0"`
	}{})
	codec.RegisterCodec("/block/proposed", struct {
		Limit   float64 `field:"limit" order:"0"`
		Page    float64 `field:"page" order:"1"`
		Address string  `field:"address" order:"2"`
	}{})
}

// Latest godoc
// @Summary      Get latest block by count
// @Description  get latest block by count
// @Tags         block
// @Accept       json
// @Produce      json
// @Success      200  {object}  bctypes.BlockResponse
// @Param        limit   query      int  true  "limit"
// @Param        page    query      int  true  "page"
// @Router       /api/block/latest [get]
func (c Block) Latest(ctx types.Context, limit, page float64) (bctypes.BlockResponse, error) {
	paginate := bctypes.NewPaginate(limit, page)
	ctx.Logger().Debug("Block.Latest", paginate)

	var blocks []bctypes.Block
	var total int64
	var err error

	var stats bctypes.Stats

	err = Cache.GetRedisCache(ctx, StatsCache, &stats)
	if err != nil {
		total = 0
	} else {
		total = stats.BlockHeight
	}

	if paginate.Offset() != 0 {
		blocks, err = bc.BlockProvider.GetLatestBlocks(ctx.DB(), paginate)
		if err != nil {
			return bctypes.BlockResponse{}, err
		}
	} else {
		err := Cache.GetRedisCache(ctx, BlockCache, &blocks)
		if err != nil {
			blocks, err = bc.BlockProvider.GetLatestBlocks(ctx.DB(), paginate)
		} else {
			if paginate.Limit() < len(blocks) {
				blocks = blocks[:paginate.Limit()]
			}
		}
	}

	return bctypes.BlockResponse{
		Blocks: blocks,
		Total:  total,
	}, nil
}

// Height godoc
// @Summary      Get block by height
// @Description  get block by height
// @Tags         block
// @Accept       json
// @Produce      json
// @Param        height   query      int  true  "height"
// @Success      200  {object}  bctypes.Block
// @Router       /api/block/height [get]
func (c Block) Height(ctx types.Context, height float64) (bctypes.Block, error) {
	h := int64(math.Round(height))
	ctx.Logger().Debug("Block.Height", h)

	return bc.BlockProvider.GetBlockByHeight(ctx.DB(), h)
}

// Proposed godoc
// @Summary      Get block by proposer
// @Description  get block by proposer
// @Tags         block
// @Accept       json
// @Produce      json
// @Success      200  {object}  bctypes.BlockResponse
// @Param        limit     query      int     true  "limit"
// @Param        page      query      int     true  "page"
// @Param        address   query      string  true  "address"
// @Router       /api/block/proposed [get]
func (c Block) Proposed(ctx types.Context, limit, page float64, address string) (bctypes.BlockResponse, error) {
	paginate := bctypes.NewPaginate(limit, page)
	ctx.Logger().Debug("Block.Proposed", paginate)

	blocks, err := bc.BlockProvider.GetProposedBlock(ctx.DB(), paginate, address)
	if err != nil {
		return bctypes.BlockResponse{}, err
	}

	total, err := bc.BlockProvider.GetProposedBlockCount(ctx.DB(), address)
	if err != nil {
		return bctypes.BlockResponse{}, err
	}

	return bctypes.BlockResponse{
		Blocks: blocks,
		Total:  total,
	}, nil
}
