package controller

import (
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/hessegg/nikto-explorer/server/api/event/codec"
	"github.com/hessegg/nikto-explorer/server/bc"
	bctypes "github.com/hessegg/nikto-explorer/server/bc/types"
	types "github.com/hessegg/nikto-explorer/types"
)

type Validator struct{}

func init() {
	codec.RegisterCodec("/validator/list", struct {
		Limit float64 `field:"limit" order:"0"`
		Page  float64 `field:"page" order:"1"`
	}{})
	codec.RegisterCodec("/validator/address", struct {
		Address string `field:"address" order:"0"`
	}{})
	codec.RegisterCodec("/validator/delegation", struct {
		Address string `field:"address" order:"0"`
		Limit   int64  `field:"limit" order:"1"`
		Page    int64  `field:"page" order:"2"`
	}{})
	codec.RegisterCodec("/validator/uptime", struct {
		Address string `field:"address" order:"0"`
	}{})
}

// List godoc
// @Summary      Get validators list
// @Description  get validators list
// @Tags         validator
// @Accept       json
// @Produce      json
// @Success      200  {object}  bctypes.ValidatorResponse
// @Param        limit     query      int     true  "limit"
// @Param        page      query      int     true  "page"
// @Router       /api/validator/list [get]
func (c Validator) List(ctx types.Context, limit, page float64) (bctypes.ValidatorResponse, error) {
	paginate := bctypes.NewPaginate(limit, page)

	var validators []bctypes.Validator
	var total int64
	var err error
	var stats bctypes.Stats
	err = Cache.GetRedisCache(ctx, StatsCache, &stats)
	if err != nil {
		total = 0
	} else {
		total = stats.ValidatorTotal
	}

	if paginate.Offset() != 0 {
		validators, err = bc.ValidatorProvider.GetValidators(ctx.DB(), paginate)
		if err != nil {
			return bctypes.ValidatorResponse{}, err
		}
	} else {
		err := Cache.GetRedisCache(ctx, ValidatorCache, &validators)
		if err != nil {
			validators, err = bc.ValidatorProvider.GetValidators(ctx.DB(), paginate)
		} else {
			if paginate.Limit() < len(validators) {
				validators = validators[:paginate.Limit()]
			}
		}
	}

	return bctypes.ValidatorResponse{
		Validators: validators,
		Total:      total,
	}, nil
}

// Address godoc
// @Summary      Get validator by address
// @Description  get validator by address
// @Tags         validator
// @Accept       json
// @Produce      json
// @Success      200  {object}  bctypes.Validator
// @Param        address   query      string  true  "address"
// @Router       /api/validator/list [get]
func (c Validator) Address(ctx types.Context, address string) (bctypes.Validator, error) {
	return bc.ValidatorProvider.GetValidator(ctx.DB(), address)
}

// Delegation godoc
// @Summary      Get validator delegations list
// @Description  get validator delegations list
// @Tags         validator
// @Accept       json
// @Produce      json
// @Success      200  {object}  bctypes.DelegationResponse
// @Param        address   query      string  true  "address"
// @Param        limit     query      int     true  "limit"
// @Param        page      query      int     true  "page"
// @Router       /api/validator/list [get]
func (c Validator) Delegation(ctx types.Context, address string, limit, page float64) (bctypes.DelegationResponse, error) {
	paginate := bctypes.NewPaginate(limit, page)
	stakingClient := staking.NewQueryClient(ctx.GRPCConn())

	res, err := stakingClient.ValidatorDelegations(ctx.Context(), &staking.QueryValidatorDelegationsRequest{
		ValidatorAddr: address,
		Pagination:    paginate.PageRequest(),
	})

	if err != nil {
		ctx.Logger().Error(err)
		return bctypes.DelegationResponse{}, err
	}

	return bctypes.DelegationResponse{
		Delegation: res.DelegationResponses,
		Total:      int64(res.Pagination.Total),
	}, nil
}

// Uptime godoc
// @Summary      Get uptime by address
// @Description  get uptime by address
// @Tags         validator
// @Accept       json
// @Produce      json
// @Success      200  {object}  bctypes.ValidatorUptime
// @Param        address   query      string  true  "address"
// @Router       /api/validator/list [get]
func (c Validator) Uptime(ctx types.Context, address string) (bctypes.ValidatorUptime, error) {
	return bc.ValidatorProvider.GetValidatorUptime(ctx.GRPCConn(), ctx.DB(), address)
}
