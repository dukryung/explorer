package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	name "github.com/hessegg/nikto/x/name/types"
)

type Account struct {
	Address       string           `json:"address"`
	Name          name.AddressName `json:"name"`
	Balance       Token            `json:"balance"`
	Balances      []Token          `json:"balances"`
	Delegations   []Delegation     `json:"delegation"`
	ReDelegations []Redelegation   `json:"redelegation"`
	UnBonding     []Unbonding      `json:"unbonding"`
}

type Delegation struct {
	Delegation types.Delegation `json:"delegation"`
	Balance    sdk.Coin         `json:"balance"`
	Moniker    string           `json:"moniker"`
}

type Redelegation struct {
	Redelegation types.Redelegation                `json:"redelegation"`
	Entries      []types.RedelegationEntryResponse `json:"entries"`
	MonikerSrc   string                            `json:"moniker_src"`
	MonikerDst   string                            `json:"moniker_dst"`
}

type Unbonding struct {
	Unbonding types.UnbondingDelegation `json:"unbonding"`
	Moniker   string                    `json:"moniker"`
}
