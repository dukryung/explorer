package types

import (
	"database/sql"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

type Validator struct {
	Detail      staking.Validator `json:"detail"`
	Moniker     string            `json:"moniker"`
	ConsPubKey  []byte            `json:"cons_pub_key"`
	ConsAddress sql.NullString    `json:"cons_address"`
	ValAddress  string            `json:"val_address"`
	Rank        int64             `json:"rank"`
}

type ValidatorUptime struct {
	LatestHeight int64         `json:"latest_height"`
	Uptime       int64         `json:"uptime"`
	Blocks       []UptimeBlock `json:"blocks"`
}

type UptimeBlock struct {
	Height int64 `json:"height"`
}

type ValidatorResponse struct {
	Validators []Validator `json:"validators"`
	Total      int64       `json:"total"`
}

type DelegationResponse struct {
	Delegation staking.DelegationResponses `json:"delegation"`
	Total      int64                       `json:"total"`
}

func (v Validator) Ed25519Address() []byte {
	return ed25519.PubKey(v.ConsPubKey[2:]).Address().Bytes()
}

type Validators struct {
	validators map[string]Validator
}

func NewValidators(val ...Validator) Validators {
	validators := Validators{}
	validators.validators = make(map[string]Validator)
	for _, validator := range val {
		validators.Add(validator)
	}
	return validators
}

func (v *Validators) Add(validator Validator) {
	v.validators[validator.ValAddress] = validator
}

func (v *Validators) GetValidators() []Validator {
	var validators []Validator

	for _, validator := range v.validators {
		validators = append(validators, validator)
	}

	return validators
}

func (v *Validators) GetValidatorByAddress(address string) Validator {
	validator, ok := v.validators[address]
	if !ok {
		return Validator{}
	}

	return validator
}
