package service

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/hessegg/nikto-explorer/server/bc"
	bctypes "github.com/hessegg/nikto-explorer/server/bc/types"
	"github.com/hessegg/nikto-explorer/types"
	"github.com/hessegg/nikto-explorer/types/config"
	"github.com/hessegg/nikto-explorer/types/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

type BlockService struct {
	ctx    types.Context
	config config.SyncConfig

	tmClient tmservice.ServiceClient
	mutex    sync.Mutex
	//dependencies
	txService        *TxService
	validatorService *ValidatorService

	currentHeight int64
	currentBlock  *tmproto.Block
}

var _ types.Service = &BlockService{}

const (
	MaxSyncBlockCount  = 500
	MaxSyncBlockWorker = 100
)

func NewBlockService(ctx types.Context, txService *TxService, validatorService *ValidatorService, config config.SyncConfig) *BlockService {
	s := &BlockService{
		ctx:              ctx.WithLogger(log.NewLogger("sync/block", config.Log)),
		config:           config,
		tmClient:         tmservice.NewServiceClient(ctx.GRPCConn()),
		mutex:            sync.Mutex{},
		txService:        txService,
		validatorService: validatorService,
		currentHeight:    1,
		currentBlock:     &tmproto.Block{},
	}

	err := s.loadLatestBlock()
	if err != nil {
		s.ctx.Logger().Error(err)
	}

	return s
}

func (s *BlockService) Run() error {
	if err := s.syncUpdatedBlock(); err != nil {
		return err
	}
	return nil
}

func (s *BlockService) syncUpdatedBlock() error {
	// if mutex locked, return
	if reflect.ValueOf(&s.mutex).Elem().FieldByName("state").Int()&1 == 1 {
		s.ctx.Logger().Debug("syncing blocks...")
		return nil
	}

	s.ctx.Logger().Debug("sync blocks...")

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.config.FastSync.Enable {
		err := s.fastSync()
		if err != nil {
			s.ctx.Logger().Error(err)
			return err
		}
	}

	s.ctx.Logger().Debug("wait for new block...")
	return nil
}

func (s *BlockService) fastSync() error {
	res, err := s.tmClient.GetLatestBlock(context.Background(), &tmservice.GetLatestBlockRequest{})
	if err != nil {
		return err
	}

	latestHeight := res.Block.Header.Height
	workerChannel := make(chan error)

	for syncCount := int64(0); s.currentHeight <= latestHeight && syncCount < MaxSyncBlockCount; {
		// get current synced height
		syncedHeight := s.currentHeight

		// set target worker count
		worker := s.targetWorker(latestHeight, syncedHeight, syncCount)
		if worker == 0 {
			break
		}

		tx, err := s.ctx.DB().Begin()
		if err != nil {
			return err
		}

		// run goroutines
		for i := int64(1); i <= worker; i++ {
			go func(height int64, result chan error) {
				err := s.syncBlockHeight(tx, height)
				if err != nil {
					result <- err
				}
				result <- nil
			}(syncedHeight+i, workerChannel)
		}

		// wait responses
		for i := int64(0); i < worker; i++ {
			select {
			// context canceled
			case <-s.ctx.Context().Done():
				s.ctx.Logger().Error(s.ctx.Context().Err())
				// rollback db
				if err := tx.Rollback(); err != nil {
					return err
				}
				return nil

			// received response from goroutine
			case err := <-workerChannel:
				if err != nil {
					// error is not nil, rollback db
					if err := tx.Rollback(); err != nil {
						return err
					}
					return err
				} else {
					// sync success, increase syncCount
					syncCount++
				}
			}
		}

		// successfully finished sync, commit db transaction
		if err := tx.Commit(); err != nil {
			return err
		}

		s.ctx.Logger().Info(fmt.Sprintf("synced block height: %v ~ %v", syncedHeight+1, syncedHeight+worker))
	}
	return nil
}

func (s *BlockService) loadLatestBlock() error {
	s.ctx.Logger().Debug("load latest block")
	var height int64
	var raw []byte

	query := `SELECT height, raw FROM block ORDER BY height DESC LIMIT 1`
	err := s.ctx.DB().QueryRow(query).Scan(&height, &raw)
	if err != nil {
		return err
	}

	s.currentHeight = height
	err = s.currentBlock.Unmarshal(raw)
	if err != nil {
		return err
	}

	s.ctx.Logger().Info(fmt.Sprintf("loaded height : %v, block : %v", s.currentHeight, s.currentBlock))
	return nil
}

func (s *BlockService) syncBlockHeight(tx *sql.Tx, height int64) error {
	//s.ctx.Logger().Debug(fmt.Sprintf("found new block height %v", height))

	// get sync target block
	res, err := s.tmClient.GetBlockByHeight(s.ctx.Context(), &tmservice.GetBlockByHeightRequest{
		Height: height,
	})

	if err != nil {
		s.ctx.Logger().Error(err)
		return err
	}
	syncBlock := res.Block

	if height-1 == 0 {
		res = &tmservice.GetBlockByHeightResponse{
			Block: &tmproto.Block{},
		}
	} else {
		res, err = s.tmClient.GetBlockByHeight(s.ctx.Context(), &tmservice.GetBlockByHeightRequest{
			Height: height - 1,
		})
		if err != nil {
			s.ctx.Logger().Error(err)
			return err
		}
	}
	prevBlock := res.Block

	err = s.txService.SyncTxs(tx, syncBlock.Data.Txs)
	if err != nil {
		s.ctx.Logger().Error(err)
		return err
	}

	validatorSet, err := bc.ValidatorProvider.GetValidatorSetByHeight(s.ctx.GRPCConn(), s.ctx.DB(), height)
	var proposer bctypes.Validator
	for _, v := range validatorSet {
		if bytes.Equal(v.Ed25519Address(), syncBlock.Header.ProposerAddress) {
			proposer = v
			break
		}
	}

	if height != 1 {
		err = s.validatorService.SyncValidatorSet(tx, syncBlock.LastCommit.Signatures, height)
		if err != nil {
			s.ctx.Logger().Error(err)
			return err
		}
	}

	raw, err := syncBlock.Marshal()
	if err != nil {
		s.ctx.Logger().Error(err)
		return err
	}

	validatorSetRaw, err := json.Marshal(validatorSet)
	if err != nil {
		s.ctx.Logger().Error(err)
		return err
	}

	var diffTime time.Duration
	if !syncBlock.Header.Time.IsZero() && !prevBlock.Header.Time.IsZero() {
		diffTime = syncBlock.Header.Time.Sub(prevBlock.Header.Time)
	}

	query := `INSERT INTO block (height, raw, cons_pubkey, block_time, diff_time, validator_set) VALUES($1, $2, $3, $4, $5, $6)`
	_, err = tx.Exec(query, syncBlock.Header.Height, raw, proposer.ConsPubKey, syncBlock.Header.Time, diffTime.Milliseconds(), validatorSetRaw)

	if err != nil {
		s.ctx.Logger().Error(err)
		return err
	}

	if s.currentHeight < syncBlock.Header.Height {
		s.currentBlock = syncBlock
		s.currentHeight = s.currentBlock.Header.Height
	}

	return nil
}

func (s *BlockService) fastSyncWorker() int64 {
	if s.config.FastSync.Worker > MaxSyncBlockWorker {
		return MaxSyncBlockWorker
	}

	return s.config.FastSync.Worker
}

func (s *BlockService) targetWorker(latestHeight, syncedHeight, syncCount int64) int64 {
	worker := int64(0)

	if diff := latestHeight - syncedHeight; diff >= s.fastSyncWorker() {
		worker = s.fastSyncWorker()
	} else {
		worker = diff
	}

	if syncCount+worker > MaxSyncBlockCount {
		worker = MaxSyncBlockCount - int64(syncCount)
	}

	return worker
}
