package bc

import (
	"database/sql"

	"github.com/hessegg/nikto-explorer/server/bc/types"
)

type block struct{}

var BlockProvider block

func (bc block) GetLatestBlocks(db *sql.DB, paginate types.Paginate) ([]types.Block, error) {
	query := `SELECT block.raw ,moniker,cons_address,val_address,diff_time FROM block block
			  LEFT JOIN validator validator
			  ON block.cons_pubkey = validator.cons_pubkey
			  ORDER BY height DESC
			  OFFSET $1
			  LIMIT $2`

	rows, err := db.Query(query, paginate.Offset(), paginate.Limit())
	if err != nil {
		return nil, err
	}

	var blocks []types.Block

	for rows.Next() {
		var block types.Block
		var raw []byte

		err = rows.Scan(&raw, &block.Moniker, &block.ConsensusAddress, &block.ValidatorAddress, &block.DiffTime)
		if err != nil {
			return nil, err
		}

		err := block.TmBlock.Unmarshal(raw)
		if err != nil {
			return nil, err
		}

		blocks = append(blocks, block)
	}

	return blocks, nil
}

func (bc block) GetBlockByHeight(db *sql.DB, height int64) (types.Block, error) {
	query := `SELECT block.raw,moniker,cons_address,val_address,diff_time FROM block block
			  LEFT JOIN validator validator
			  ON block.cons_pubkey = validator.cons_pubkey
			  WHERE height = $1
`

	var block types.Block
	var raw []byte
	err := db.QueryRow(query, height).Scan(&raw, &block.Moniker, &block.ConsensusAddress, &block.ValidatorAddress, &block.DiffTime)
	if err != nil {
		return types.Block{}, err
	}

	err = block.TmBlock.Unmarshal(raw)
	if err != nil {
		return types.Block{}, err
	}

	return block, nil
}

func (bc block) GetProposedBlock(db *sql.DB, paginate types.Paginate, address string) ([]types.Block, error) {
	query := `SELECT block.raw,moniker,cons_address,val_address,diff_time FROM block block
			  LEFT JOIN validator validator
			  ON block.cons_pubkey = validator.cons_pubkey
			  WHERE val_address = $1
			  ORDER BY height DESC
			  OFFSET $2
			  LIMIT $3`

	rows, err := db.Query(query, address, paginate.Offset(), paginate.Limit())
	if err != nil {
		return nil, err
	}

	var blocks []types.Block

	for rows.Next() {
		var block types.Block
		var raw []byte

		err = rows.Scan(&raw, &block.Moniker, &block.ConsensusAddress, &block.ValidatorAddress, &block.DiffTime)
		if err != nil {
			return nil, err
		}

		err := block.TmBlock.Unmarshal(raw)
		if err != nil {
			return nil, err
		}

		blocks = append(blocks, block)
	}

	return blocks, nil
}

func (bc block) GetProposedBlockCount(db *sql.DB, address string) (int64, error) {
	query := `SELECT count(*) FROM block
			  LEFT JOIN validator
			  ON block.cons_pubkey = validator.cons_pubkey
			  WHERE val_address = $1`

	var count int64
	err := db.QueryRow(query, address).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
