package controller

import (
	"github.com/hessegg/nikto-explorer/server/api/event/codec"
	"github.com/hessegg/nikto-explorer/server/bc"
	bctypes "github.com/hessegg/nikto-explorer/server/bc/types"
	"github.com/hessegg/nikto-explorer/types"
)

type Token struct{}

func init() {
	codec.RegisterCodec("/token/list", struct {
		Limit float64 `field:"limit" order:"0"`
		Page  float64 `field:"page" order:"1"`
	}{})
	codec.RegisterCodec("/token/by_denom", struct {
		Denom string `field:"denom" order:"0"`
	}{})

}

// List godoc
// @Summary      Get latest block by count
// @Description  get latest block by count
// @Tags         token
// @Accept       json
// @Produce      json
// @Param        limit   query      int  true  "limit"
// @Param        page    query      int  true  "page"
// @Router       /api/token/list [get]
func (c Token) List(ctx types.Context, limit, page float64) (bctypes.TokenResponse, error) {
	paginate := bctypes.NewPaginate(limit, page)
	ctx.Logger().Debug("Token.List", paginate)

	var tokens []bctypes.Token
	var total int64
	var err error
	var stats bctypes.Stats

	err = Cache.GetRedisCache(ctx, StatsCache, &stats)
	if err != nil {
		total = 0
	} else {
		total = stats.TokenTotal
	}

	if paginate.Offset() != 0 {
		tokens, err = bc.TokenProvider.GetTokenList(ctx.DB(), paginate)
		if err != nil {
			return bctypes.TokenResponse{}, err
		}
	} else {
		err := Cache.GetRedisCache(ctx, TokenCache, &tokens)
		if err != nil {
			tokens, err = bc.TokenProvider.GetTokenList(ctx.DB(), paginate)
		} else {
			if paginate.Limit() < len(tokens) {
				tokens = tokens[:paginate.Limit()]
			}
		}
	}

	return bctypes.TokenResponse{
		Tokens: tokens,
		Total:  total,
	}, nil
}

// ByDenom godoc
// @Summary      Get token by denominator
// @Description  get token by denominator
// @Tags         token
// @Accept       json
// @Produce      json
// @Param        denom   query      string  true  "denom"
// @Router       /api/token/bydenom [get]
func (c Token) ByDenom(ctx types.Context, denom string) (bctypes.Token, error) {
	return bc.TokenProvider.GetToken(ctx.DB(), denom)
}
