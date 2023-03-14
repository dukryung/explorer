package service

import (
	"database/sql"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	servicetypes "github.com/hessegg/nikto-explorer/server/sync/service/types"
	"github.com/hessegg/nikto-explorer/types"
	"github.com/hessegg/nikto-explorer/types/config"
	"github.com/hessegg/nikto-explorer/types/log"
)

type TxService struct {
	ctx           types.Context
	tmClient      tmservice.ServiceClient
	txClient      tx.ServiceClient
	txDecoder     sdk.TxDecoder
	txJSONEncoder sdk.TxEncoder

	config config.SyncConfig
	mutex  sync.Mutex
}

var _ types.Service = &TxService{}

const (
	MaxSyncTxCount  = 500
	MaxSyncTxWorker = 100
)

func NewTxService(ctx types.Context, config config.SyncConfig) *TxService {
	service := &TxService{
		ctx:           ctx.WithLogger(log.NewLogger("sync/TxService", config.Log)),
		tmClient:      tmservice.NewServiceClient(ctx.GRPCConn()),
		txClient:      tx.NewServiceClient(ctx.GRPCConn()),
		txDecoder:     ctx.TxConfig().TxDecoder(),
		txJSONEncoder: ctx.TxConfig().TxJSONEncoder(),
		config:        config,
		mutex:         sync.Mutex{},
	}

	return service
}

func (s *TxService) Run() error {
	if err := s.syncTxs(); err != nil {
		return err
	}
	return nil
}

// SyncTxs sync txs round 1
// called by block service
func (s *TxService) SyncTxs(tx *sql.Tx, txs [][]byte) error {
	if len(txs) == 0 {
		return nil
	}

	s.ctx.Logger().Info(fmt.Sprintf("found %v txs", len(txs)))

	for _, raw := range txs {
		txHash := types.NewTxHash(raw)
		s.ctx.Logger().Debug(fmt.Sprintf("synced tx %v", txHash.String()))

		query := `INSERT INTO transaction (txhash, txbody) VALUES ($1, $2)`
		_, err := tx.Exec(query, txHash.String(), raw)
		if err != nil {
			s.ctx.Logger().Error(err)
			return err
		}
	}

	return nil
}

// syncTxs sync txs round 2
// Update the registered tx data.
func (s *TxService) syncTxs() error {
	// if mutex locked, return
	if reflect.ValueOf(&s.mutex).Elem().FieldByName("state").Int()&1 == 1 {
		s.ctx.Logger().Debug("syncing txs...")
		return nil
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	db, err := s.ctx.DB().Begin()
	if err != nil {
		s.ctx.Logger().Error(err)
		return err
	}

	var pending int64
	query := fmt.Sprintf(`SELECT count(*) FROM transaction WHERE updated = FALSE`)
	err = db.QueryRow(query).Scan(&pending)
	if err != nil {
		s.ctx.Logger().Error(err)
		return err
	}

	if pending == 0 {
		if err := db.Rollback(); err != nil {
			return err
		}
		return nil
	}

	var syncCount int64
	for syncCount = int64(0); syncCount < MaxSyncTxCount && syncCount < pending; {
		workerCount := s.targetWorker(pending, syncCount)

		query = fmt.Sprintf(`SELECT txhash FROM transaction WHERE updated = FALSE LIMIT %v`, workerCount)
		rows, err := db.Query(query)
		if err != nil {
			s.ctx.Logger().Error(err)
			return err
		}

		txHashes := make([]string, 0)
		for rows.Next() {
			var txHash string
			err := rows.Scan(&txHash)
			if err != nil {
				// error is not nil, rollback db
				if err := db.Rollback(); err != nil {
					return err
				}
				return err
			}
			txHashes = append(txHashes, txHash)
		}

		workerChannel := make(chan error)
		for i := int64(0); i < workerCount; i++ {
			go func(txHash string, result chan error) {
				err := s.syncTxHash(db, txHash)
				if err != nil {
					result <- err
				}
				result <- nil
			}(txHashes[i], workerChannel)
		}

		for i := int64(0); i < workerCount; i++ {
			select {
			// context canceled
			case <-s.ctx.Context().Done():
				s.ctx.Logger().Error(s.ctx.Context().Err())
				// rollback db
				if err := db.Rollback(); err != nil {
					return err
				}
				return nil

			// received response from goroutine
			case err := <-workerChannel:
				if err != nil {
					// error is not nil, rollback db
					if err := db.Rollback(); err != nil {
						return err
					}
					return err
				} else {
					syncCount++
				}
			}
		}
	}

	// successfully finished sync, commit db transaction
	if err := db.Commit(); err != nil {
		s.ctx.Logger().Error(err)
		return err
	}

	s.ctx.Logger().Info("updated txs.\tpending : ", pending-syncCount)

	return nil
}

func (s *TxService) syncTxHash(db *sql.Tx, txHash string) error {
	res, err := s.txClient.GetTx(s.ctx.Context(), &tx.GetTxRequest{
		Hash: txHash,
	})

	if err != nil {
		return err
	}

	txResponse := res.TxResponse
	txTime, err := time.Parse(time.RFC3339, txResponse.Timestamp)
	if err != nil {
		return err
	}

	decodedTx, err := s.txDecoder(txResponse.Tx.Value)
	if err != nil {
		return err
	}

	encodedTx, err := s.txJSONEncoder(decodedTx)
	if err != nil {
		return err
	}

	address, err := servicetypes.ParseSignerAddress(encodedTx)
	if err != nil {
		return err
	}

	query := `UPDATE transaction 
				  SET  height = $1, txbody = $2, timestamp = $3, updated = TRUE, 
				       code = $4, sender = $5
				  WHERE txhash = $6`
	_, err = db.Exec(query, txResponse.Height, txResponse.Tx.Value, txTime, txResponse.Code, address.String(), txHash)
	if err != nil {
		return err
	}

	s.ctx.Logger().Debug(fmt.Sprintf("updated tx %v", txHash))
	return nil
}

func (s *TxService) fastSyncWorker() int64 {
	if s.config.FastSync.Worker > MaxSyncTxWorker {
		return MaxSyncTxWorker
	}

	return s.config.FastSync.Worker
}

func (s *TxService) targetWorker(pending, syncCount int64) int64 {
	workerCount := s.fastSyncWorker()
	if workerCount+syncCount > pending {
		workerCount = pending - syncCount
	}

	if workerCount+syncCount > MaxSyncTxCount {
		workerCount = MaxSyncTxCount - syncCount
	}

	return workerCount
}
