package controller

import (
	"github.com/hessegg/nikto-explorer/server/api/event/codec"
	"github.com/hessegg/nikto-explorer/server/bc"
	bctypes "github.com/hessegg/nikto-explorer/server/bc/types"
	"github.com/hessegg/nikto-explorer/types"
)

type Stats struct{}

func init() {
	codec.RegisterCodec("/stats/now", struct{}{})
	codec.RegisterCodec("/stats/tx", struct{}{})
}

// Now godoc
// @Summary      Get stats of explorer
// @Description  get stats of explorer
// @Tags         stats
// @Accept       json
// @Produce      json
// @Success      200  {object}  bctypes.Stats
// @Router       /api/stats/now [get]
func (c Stats) Now(ctx types.Context) (bctypes.Stats, error) {
	ctx.Logger().Debug("Stats.Now")

	var stats bctypes.Stats
	err := Cache.GetRedisCache(ctx, StatsCache, &stats)
	if err != nil {
		return bc.StatsProvider.GetStats(ctx.DB())
	}

	return stats, nil
}

// Tx godoc
// @Summary      Get stats of explorer
// @Description  get stats of explorer
// @Tags         stats
// @Accept       json
// @Produce      json
// @Success      200  {object}  []bctypes.TxStats
// @Router       /api/stats/tx [get]
func (c Stats) Tx(ctx types.Context) ([]bctypes.TxStats, error) {
	ctx.Logger().Debug("Stats.Tx")

	var txStats []bctypes.TxStats
	err := Cache.GetRedisCache(ctx, TxStatsCache, &txStats)
	if err != nil {
		return bc.StatsProvider.GetTxStats(ctx.DB())
	}

	return txStats, nil
}
