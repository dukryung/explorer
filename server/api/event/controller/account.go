package controller

import (
	"context"

	"github.com/cosmos/cosmos-sdk/types/query"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/hessegg/nikto-explorer/server/api/event/codec"
	"github.com/hessegg/nikto-explorer/server/bc"
	bctypes "github.com/hessegg/nikto-explorer/server/bc/types"
	"github.com/hessegg/nikto-explorer/types"
	name "github.com/hessegg/nikto/x/name/types"
)

type Account struct{}

func init() {
	codec.RegisterCodec("/account/address", struct {
		Address string `field:"address" order:"0"`
	}{})
	codec.RegisterCodec("/account/balance", struct {
		Address string `field:"address" order:"0"`
		Denom   string `field:"denom" order:"1"`
	}{})
	codec.RegisterCodec("/account/allbalances", struct {
		Address     string `field:"address" order:"0"`
		BondedDenom string `field:"bondedDenom" order:"1"`
	}{})
	codec.RegisterCodec("/account/delegation", struct {
		Address string `field:"address" order:"0"`
	}{})
	codec.RegisterCodec("/account/redelegation", struct {
		Address string `field:"address" order:"0"`
	}{})
	codec.RegisterCodec("/account/unbonding", struct {
		Address string `field:"address" order:"0"`
	}{})
}

// Address godoc
// @Summary      Get account by address
// @Description  Get account by address
// @Tags         account
// @Accept       json
// @Produce      json
// @Success      200  {object}  bctypes.Account
// @Param        address   query      string  true  "address"
// @Router       /api/account/address [get]
func (a Account) Address(ctx types.Context, address string) (bctypes.Account, error) {
	account := bctypes.Account{}
	account.Address = address

	var err error

	account.Delegations, err = a.Delegation(ctx, address)
	if err != nil {
		return account, err
	}

	account.UnBonding, err = a.Unbonding(ctx, address)
	if err != nil {
		return account, err
	}

	account.ReDelegations, err = a.Redelegation(ctx, address)
	if err != nil {
		return account, err
	}

	account.Name, err = a.Name(ctx, address)
	if err != nil {
		return account, err
	}

	ctx.Logger().Info("Account.Address", address)
	return account, nil
}

// Balance godoc
// @Summary      Get balance by denom
// @Description  Get balance by denom
// @Tags         account
// @Accept       json
// @Produce      json
// @Success      200  {object}  bctypes.Token
// @Param        address   query      string  true  "address"
// @Param        denom   query      string  true  "denom"
// @Router       /api/account/balance [get]
func (a Account) Balance(ctx types.Context, address, denom string) (bctypes.Token, error) {
	bankClient := bank.NewQueryClient(ctx.GRPCConn())
	balance, err := bankClient.Balance(ctx.Context(), &bank.QueryBalanceRequest{
		Address: address,
		Denom:   denom,
	})
	if err != nil {
		return bctypes.Token{}, err
	}

	token, _ := bc.TokenProvider.GetToken(ctx.DB(), denom)
	token.Coin = balance.Balance

	ctx.Logger().Debug("Account.Balance", token)
	return token, nil
}

// AllBalances godoc
// @Summary      Get all balances by address
// @Description  Get all balances by address
// @Tags         account
// @Accept       json
// @Produce      json
// @Success      200  {object}  []bctypes.Token
// @Param        address   query      string  true  "address"
// @Param        bondedDenom   query      string  true  "bondedDenom"
// @Router       /api/account/address [get]
func (a Account) AllBalances(ctx types.Context, address, bondedDenom string) ([]bctypes.Token, error) {
	bankClient := bank.NewQueryClient(ctx.GRPCConn())
	balance, err := bankClient.AllBalances(context.Background(), &bank.QueryAllBalancesRequest{
		Address: address,
		Pagination: &query.PageRequest{
			Limit:  100,
			Offset: 0,
		},
	})

	if err != nil {
		return nil, err
	}

	var denoms []string
	for _, token := range balance.Balances {
		denoms = append(denoms, token.Denom)
	}

	tokens, err := bc.TokenProvider.GetTokens(ctx.DB(), denoms)
	if err != nil {
		return nil, err
	}

	for _, token := range tokens {
		token.Coin.Amount = balance.Balances.AmountOf(token.Coin.Denom)
	}

	// No need. AllBalances returns platform coin also
	// netToken, err := a.Balance(ctx, address, bondedDenom)
	// if err != nil {
	// 	return nil, err
	// }

	// tokens = append([]bctypes.Token{netToken}, tokens...)

	return tokens, nil
}

// Delegation godoc
// @Summary      Get delegations by address
// @Description  Get delegations by address
// @Tags         account
// @Accept       json
// @Produce      json
// @Success      200  {object}  []bctypes.Delegation
// @Param        address   query      string  true  "address"
// @Router       /api/account/delegation [get]
func (a Account) Delegation(ctx types.Context, address string) ([]bctypes.Delegation, error) {
	stakingClient := staking.NewQueryClient(ctx.GRPCConn())
	response, err := stakingClient.DelegatorDelegations(ctx.Context(), &staking.QueryDelegatorDelegationsRequest{
		DelegatorAddr: address,
	})

	if err != nil {
		return nil, err
	}

	val, err := bc.ValidatorProvider.GetValidators(ctx.DB(), bctypes.NewPaginate(1000, 0))
	if err != nil {
		return nil, err
	}

	validators := bctypes.NewValidators(val...)
	var delegations []bctypes.Delegation

	for _, resp := range response.DelegationResponses {
		var delegation bctypes.Delegation

		delegation.Delegation = resp.Delegation
		delegation.Balance = resp.Balance
		delegation.Moniker = validators.GetValidatorByAddress(resp.Delegation.ValidatorAddress).Moniker

		delegations = append(delegations, delegation)
	}

	return delegations, nil
}

// Redelegation godoc
// @Summary      Get redelegations by address
// @Description  Get redelegations by address
// @Tags         account
// @Accept       json
// @Produce      json
// @Success      200  {object}  []bctypes.Redelegation
// @Param        address   query      string  true  "address"
// @Router       /api/account/redelegation [get]
func (a Account) Redelegation(ctx types.Context, address string) ([]bctypes.Redelegation, error) {
	stakingClient := staking.NewQueryClient(ctx.GRPCConn())
	response, err := stakingClient.Redelegations(ctx.Context(), &staking.QueryRedelegationsRequest{
		DelegatorAddr: address,
	})

	if err != nil {
		return nil, err
	}

	val, err := bc.ValidatorProvider.GetValidators(ctx.DB(), bctypes.NewPaginate(1000, 0))
	if err != nil {
		return nil, err
	}

	validators := bctypes.NewValidators(val...)
	var redelegations []bctypes.Redelegation

	for _, resp := range response.RedelegationResponses {
		var redelegation bctypes.Redelegation

		redelegation.Redelegation = resp.Redelegation
		redelegation.Entries = resp.Entries
		redelegation.MonikerSrc = validators.GetValidatorByAddress(resp.Redelegation.ValidatorSrcAddress).Moniker
		redelegation.MonikerDst = validators.GetValidatorByAddress(resp.Redelegation.ValidatorDstAddress).Moniker

		redelegations = append(redelegations, redelegation)
	}

	return redelegations, nil
}

// Unbonding godoc
// @Summary      Get unbonding by address
// @Description  Get unbonding by address
// @Tags         account
// @Accept       json
// @Produce      json
// @Success      200  {object}  []bctypes.Unbonding
// @Param        address   query      string  true  "address"
// @Router       /api/account/unbonding [get]
func (a Account) Unbonding(ctx types.Context, address string) ([]bctypes.Unbonding, error) {
	stakingClient := staking.NewQueryClient(ctx.GRPCConn())
	response, err := stakingClient.DelegatorUnbondingDelegations(ctx.Context(), &staking.QueryDelegatorUnbondingDelegationsRequest{
		DelegatorAddr: address,
		Pagination: &query.PageRequest{
			CountTotal: true,
		},
	})

	if err != nil {
		return nil, err
	}

	val, err := bc.ValidatorProvider.GetValidators(ctx.DB(), bctypes.NewPaginate(1000, 0))
	if err != nil {
		return nil, err
	}

	validators := bctypes.NewValidators(val...)
	var unbondings []bctypes.Unbonding

	for _, resp := range response.UnbondingResponses {
		var unbonding bctypes.Unbonding

		unbonding.Unbonding = resp
		unbonding.Moniker = validators.GetValidatorByAddress(resp.ValidatorAddress).Moniker

		unbondings = append(unbondings, unbonding)
	}

	return unbondings, nil
}

func (a Account) Name(ctx types.Context, address string) (name.AddressName, error) {
	nameClient := name.NewQueryClient(ctx.GRPCConn())
	response, err := nameClient.Address(ctx.Context(), &name.QueryAddressRequest{
		Address: address,
	})

	if err != nil {
		return name.AddressName{}, err
	}

	return *response.AddressName, nil
}
