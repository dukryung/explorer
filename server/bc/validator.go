package bc

import (
	"bytes"
	"context"
	"database/sql"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	slashing "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/hessegg/nikto-explorer/server/bc/types"
	"google.golang.org/grpc"
)

type validator struct{}

const (
	MAX_UPTIME_BLOCK = 100
)

var ValidatorProvider validator

func (bc validator) GetValidators(db *sql.DB, paginate types.Paginate) ([]types.Validator, error) {
	query := `SELECT val_address, cons_pubkey, cons_address, moniker, raw, rank 
       		  FROM validator
			  ORDER BY rank
			  OFFSET $1 
              LIMIT $2`

	rows, err := db.Query(query, paginate.Offset(), paginate.Limit())

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var validators []types.Validator

	for rows.Next() {
		var validator types.Validator
		var raw []byte

		err := rows.Scan(
			&validator.ValAddress,
			&validator.ConsPubKey,
			&validator.ConsAddress,
			&validator.Moniker,
			&raw,
			&validator.Rank,
		)
		if err != nil {
			return nil, err
		}

		err = validator.Detail.Unmarshal(raw)
		validator.Detail.ConsensusPubkey = nil
		if err != nil {
			return nil, err
		}

		validators = append(validators, validator)
	}

	return validators, nil
}

func (bc validator) GetValidator(db *sql.DB, address string) (types.Validator, error) {
	var validator types.Validator
	var raw []byte
	query := `SELECT val_address, cons_pubkey, cons_address, moniker, raw, rank
				  FROM validator 
				  WHERE val_address = $1`

	err := db.QueryRow(query, address).Scan(
		&validator.ValAddress,
		&validator.ConsPubKey,
		&validator.ConsAddress,
		&validator.Moniker,
		&raw,
		&validator.Rank,
	)
	if err != nil {
		return types.Validator{}, err
	}

	err = validator.Detail.Unmarshal(raw)
	if err != nil {
		return types.Validator{}, err
	}

	validator.Detail.ConsensusPubkey = nil

	return validator, nil
}

func (bc validator) GetValidatorSetByHeight(conn *grpc.ClientConn, db *sql.DB, height int64) ([]types.Validator, error) {
	tmClient := tmservice.NewServiceClient(conn)
	res, err := tmClient.GetValidatorSetByHeight(context.Background(), &tmservice.GetValidatorSetByHeightRequest{
		Height: height,
	})
	if err != nil {
		return nil, err
	}

	var validators []types.Validator

	//TODO bc can be multiple
	for _, v := range res.Validators {
		var validator types.Validator
		query := `SELECT val_address, cons_pubkey, cons_address, moniker, rank
				  FROM validator 
				  WHERE cons_pubkey = $1`
		err := db.QueryRow(query, v.PubKey.Value).Scan(
			&validator.ValAddress,
			&validator.ConsPubKey,
			&validator.ConsAddress,
			&validator.Moniker,
			&validator.Rank,
		)

		if err != nil {
			return nil, err
		}

		validators = append(validators, validator)
	}

	return validators, nil
}

func (bc validator) GetProposerByHeight(conn *grpc.ClientConn, db *sql.DB, height int64) (types.Validator, error) {
	tmClient := tmservice.NewServiceClient(conn)
	res, err := tmClient.GetBlockByHeight(context.Background(), &tmservice.GetBlockByHeightRequest{
		Height: height,
	})
	if err != nil {
		return types.Validator{}, err
	}

	proposerAddress := res.Block.Header.ProposerAddress

	//get validators
	validators, err := bc.GetValidatorSetByHeight(conn, db, height)
	if err != nil {
		return types.Validator{}, err
	}

	//compare validators
	var proposer types.Validator
	for _, v := range validators {
		if bytes.Equal(v.Ed25519Address(), proposerAddress) {
			proposer = v
			break
		}
	}

	return proposer, nil
}

func (bc validator) GetValidatorUptime(conn *grpc.ClientConn, db *sql.DB, address string) (types.ValidatorUptime, error) {
	var consAddress string
	query := `SELECT cons_address from validator where val_address=$1`

	err := db.QueryRow(query, address).Scan(&consAddress)
	if err != nil {
		return types.ValidatorUptime{}, err
	}
	var validatorUptime types.ValidatorUptime
	// calculate uptime
	slashingClient := slashing.NewQueryClient(conn)

	paramRes, err := slashingClient.Params(context.Background(), &slashing.QueryParamsRequest{})
	if err != nil {
		return types.ValidatorUptime{}, err
	}

	infoRes, err := slashingClient.SigningInfo(context.Background(), &slashing.QuerySigningInfoRequest{
		ConsAddress: consAddress,
	})
	if err != nil {
		return types.ValidatorUptime{}, err
	}

	signedBlocksWindow := paramRes.Params.SignedBlocksWindow
	missedBlockCounter := infoRes.ValSigningInfo.MissedBlocksCounter

	//println(signedBlocksWindow, missedBlockCounter)

	uptime := (signedBlocksWindow - missedBlockCounter) / signedBlocksWindow * 100
	// Set Uptime
	validatorUptime.Uptime = uptime

	// Get Latest Block
	latestBlock, err := BlockProvider.GetLatestBlocks(db, types.NewPaginate(1, 0))
	if err != nil {
		return types.ValidatorUptime{}, err
	}
	latestHeight := latestBlock[0].TmBlock.Header.Height

	// Set Latest Height
	validatorUptime.LatestHeight = latestHeight

	// Get Uptime Blocks
	query = `SELECT height from uptime 
              WHERE height >= $1 
                AND height < $2
                AND cons_address = $3`

	rows, err := db.Query(query, latestHeight-MAX_UPTIME_BLOCK, latestHeight, consAddress)
	if err != nil {
		return types.ValidatorUptime{}, err
	}

	// Initialize Blocks
	validatorUptime.Blocks = make([]types.UptimeBlock, 0)
	for rows.Next() {
		var block types.UptimeBlock
		err := rows.Scan(&block.Height)
		if err != nil {
			return types.ValidatorUptime{}, err
		}

		// Set Blocks
		validatorUptime.Blocks = append(validatorUptime.Blocks, block)
	}

	return validatorUptime, nil
}
