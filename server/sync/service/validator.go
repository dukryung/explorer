package service

import (
	"context"
	"database/sql"
	"sync"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	sdk "github.com/cosmos/cosmos-sdk/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/hessegg/nikto-explorer/server/bc"
	bctypes "github.com/hessegg/nikto-explorer/server/bc/types"
	"github.com/hessegg/nikto-explorer/server/util"
	"github.com/hessegg/nikto-explorer/types"
	"github.com/hessegg/nikto-explorer/types/config"
	"github.com/hessegg/nikto-explorer/types/log"
	tmblock "github.com/tendermint/tendermint/proto/tendermint/types"
)

type ValidatorService struct {
	ctx           types.Context
	tmClient      tmservice.ServiceClient
	stakingClient staking.QueryClient
	mutex         sync.Mutex
}

var _ types.Service = &ValidatorService{}

func NewValidatorService(ctx types.Context, config config.LogConfig) *ValidatorService {
	service := &ValidatorService{
		ctx:           ctx.WithLogger(log.NewLogger("sync/validator", config)),
		tmClient:      tmservice.NewServiceClient(ctx.GRPCConn()),
		stakingClient: staking.NewQueryClient(ctx.GRPCConn()),
		mutex:         sync.Mutex{},
	}

	if err := service.syncValidator(); err != nil {
		service.ctx.Logger().Error(err)
	}

	return service
}

func (s *ValidatorService) Run() error {
	if err := s.syncValidator(); err != nil {
		return err
	}
	return nil
}

func (s *ValidatorService) syncValidator() error {
	s.ctx.Logger().Debug("sync validator set...")

	resp, err := s.stakingClient.Validators(context.Background(), &staking.QueryValidatorsRequest{})
	if err != nil {
		return err
	}

	db, err := s.ctx.DB().Begin()
	if err != nil {
		return err
	}

	for _, validator := range resp.Validators {
		raw, err := validator.Marshal()
		if err != nil {
			return err
		}

		query := `INSERT INTO validator (val_address, cons_address ,cons_pubkey, moniker, raw, tokens) 
				  VALUES ($1, $2, $3, $4, $5, $6)
				  ON CONFLICT(cons_pubkey) 
				  DO UPDATE 
				  SET moniker = $4, raw = $5, tokens = $6
`

		_, err = db.Exec(query,
			validator.OperatorAddress,
			util.ConsPubKeyToBech32(validator.ConsensusPubkey.Value[2:]),
			validator.ConsensusPubkey.Value,
			validator.Description.Moniker,
			raw,
			validator.Tokens.String())

		if err != nil {
			s.ctx.Logger().Error(err)
			if err = db.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	if err := s.updateValidatorRank(db); err != nil {
		s.ctx.Logger().Error(err)
		if err = db.Rollback(); err != nil {
			return err
		}
		return err
	}

	if err := db.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *ValidatorService) updateValidatorRank(tx *sql.Tx) error {
	query := `UPDATE validator v1
			  SET rank = v2.rank
			  FROM (
				SELECT cons_address,
	 			   	   RANK() OVER (ORDER BY CAST(tokens AS numeric) DESC, moniker) as rank
				FROM validator
			  ) v2
			  WHERE v1.cons_address = v2.cons_address;
`
	_, err := tx.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

// SyncValidatorSet
// sync when block syncs
func (s *ValidatorService) SyncValidatorSet(tx *sql.Tx, sig []tmblock.CommitSig, height int64) error {
	validators, err := bc.ValidatorProvider.GetValidators(s.ctx.DB(), bctypes.NewPaginate(1000, 0))
	if err != nil {
		s.ctx.Logger().Error(err)
		return err
	}

	diff := validatorDifference(validators, sig)
	for _, val := range diff {
		query := `INSERT INTO uptime (height, cons_address) VALUES ($1, $2)`
		_, err = tx.Exec(query, height, val)
		if err != nil {
			s.ctx.Logger().Error(err)
			return err
		}
	}

	return nil
}

func validatorDifference(val1 []bctypes.Validator, val2 []tmblock.CommitSig) []string {
	var diff []string
	m := map[string]int{}

	for _, val := range val1 {
		addr := sdk.ConsAddress(val.Ed25519Address()).String()
		m[addr] = 1
	}

	for _, val := range val2 {
		addr := sdk.ConsAddress(val.ValidatorAddress).String()
		m[addr] = m[addr] + 1
	}

	for key, val := range m {
		if val == 1 {
			diff = append(diff, key)
		}
	}

	return diff
}
