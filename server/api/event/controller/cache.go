package controller

import (
	"encoding/json"
	"sync"

	"github.com/hessegg/nikto-explorer/server/bc"
	bctypes "github.com/hessegg/nikto-explorer/server/bc/types"
	"github.com/hessegg/nikto-explorer/types"
)

type cache struct {
	blocks  []bctypes.Block
	txs     []bctypes.Tx
	stats   bctypes.Stats
	tokens  []bctypes.Token
	txStats []bctypes.TxStats

	stateMux *sync.Mutex
	state    map[int]bool
	mux      map[int]*sync.Mutex
}

const (
	BlockCache = 0 + iota
	TxCache
	StatsCache
	TokenCache
	TxStatsCache
	ValidatorCache
	NftCache
)

var (
	caches = []int{
		BlockCache,
		TxCache,
		StatsCache,
		TokenCache,
		TxStatsCache,
		ValidatorCache,
		NftCache,
	}

	redisKey = []string{
		"block",
		"tx",
		"stats",
		"token",
		"txStats",
		"validators",
		"nftokens",
	}
)

var Cache cache

func init() {
	Cache.state = make(map[int]bool)
	Cache.mux = make(map[int]*sync.Mutex)

	initCaches()
}

func initCaches() {

}

func (c *cache) GetRedisCache(ctx types.Context, cache int, object interface{}) error {
	res, err := ctx.Redis().Get(ctx.Context(), redisKey[cache]).Result()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(res), &object)
	if err != nil {
		return err
	}

	return nil
}

func (c *cache) SetBlockCache(ctx types.Context) error {
	blocks, err := bc.BlockProvider.GetLatestBlocks(ctx.DB(), bctypes.DefaultPaginate())
	if err != nil {
		ctx.Logger().Error(err)
		return err
	}

	err = c.setRedisCache(ctx, redisKey[BlockCache], blocks)
	if err != nil {
		ctx.Logger().Error(err)
	}

	if err != nil {
		ctx.Logger().Error(err)
		return err
	}
	return nil
}

func (c *cache) SetTxCache(ctx types.Context) error {
	txs, err := bc.TxProvider.GetLatestTxs(ctx.DB(), bctypes.DefaultPaginate())
	if err != nil {
		ctx.Logger().Error(err)
		return err
	}

	err = c.setRedisCache(ctx, redisKey[TxCache], txs)
	if err != nil {
		ctx.Logger().Error(err)
	}
	return nil
}

func (c *cache) SetStatsCache(ctx types.Context) error {
	stats, err := bc.StatsProvider.GetStats(ctx.DB())
	if err != nil {
		ctx.Logger().Error(err)
		return err
	}

	err = c.setRedisCache(ctx, redisKey[StatsCache], stats)
	if err != nil {
		ctx.Logger().Error(err)
	}
	return nil
}

func (c *cache) SetTxStatsCache(ctx types.Context) error {
	txStats, err := bc.StatsProvider.GetTxStats(ctx.DB())
	if err != nil {
		ctx.Logger().Error(err)
		return err
	}

	err = c.setRedisCache(ctx, redisKey[TxStatsCache], txStats)
	if err != nil {
		ctx.Logger().Error(err)
		return err
	}
	return nil
}

func (c *cache) SetTokenCache(ctx types.Context) error {
	tokenList, err := bc.TokenProvider.GetTokenList(ctx.DB(), bctypes.DefaultPaginate())
	if err != nil {
		ctx.Logger().Error(err)
		return err
	}

	err = c.setRedisCache(ctx, redisKey[TokenCache], tokenList)
	if err != nil {
		ctx.Logger().Error(err)
		return err
	}
	return nil
}

func (c *cache) SetNftCache(ctx types.Context) error {
	infoList, err := bc.NftProvider.GetNftList(ctx.DB(), bctypes.DefaultPaginate())
	if err != nil {
		ctx.Logger().Error(err)
		return err
	}

	err = c.setRedisCache(ctx, redisKey[NftCache], infoList)
	if err != nil {
		ctx.Logger().Error(err)
		return err
	}
	return nil
}

func (c *cache) SetValidatorCache(ctx types.Context) error {
	validators, err := bc.ValidatorProvider.GetValidators(ctx.DB(), bctypes.DefaultPaginate())
	if err != nil {
		ctx.Logger().Error(err)
		return err
	}

	err = c.setRedisCache(ctx, redisKey[ValidatorCache], validators)
	if err != nil {
		ctx.Logger().Error(err)
		return err
	}
	return nil
}

func (c *cache) setRedisCache(ctx types.Context, key string, value interface{}) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = ctx.Redis().
		Set(ctx.Context(), key, string(b), 0).
		Err()

	if err != nil {
		return err
	}

	return nil
}
