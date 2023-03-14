package bc

import (
	"database/sql"
	"encoding/json"

	"github.com/hessegg/nikto-explorer/server/bc/types"
	klaatoo "github.com/hessegg/nikto/app"
)

type tx struct{}

var TxProvider tx

func (bc tx) GetLatestTxs(db *sql.DB, paginate types.Paginate) ([]types.Tx, error) {
	query := `SELECT height,txbody,txhash,timestamp,code,sender
			  FROM transaction
			  WHERE updated = TRUE
			  ORDER BY idx DESC 
			  OFFSET $1
			  LIMIT $2`

	rows, err := db.Query(query, paginate.Offset(), paginate.Limit())
	if err != nil {
		return nil, err
	}

	var txs []types.Tx
	for rows.Next() {
		var tx types.Tx
		var txBody []byte

		err = rows.Scan(&tx.Height, &txBody, &tx.TxHash, &tx.TimeStamp, &tx.Code, &tx.Sender)

		if err != nil {
			return nil, err
		}

		jsonTx, err := bc.encodeTxBodyJSON(txBody)
		if err != nil {
			return nil, err
		}

		tx.Tx = jsonTx
		txs = append(txs, tx)
	}

	return txs, nil
}

func (bc tx) GetTxByHash(db *sql.DB, hash string) (types.Tx, error) {
	var tx types.Tx
	var txBody []byte

	query := `SELECT height,txbody,txhash,timestamp,code,sender 
			  FROM transaction 
			  WHERE txhash=$1`
	err := db.QueryRow(query, hash).Scan(&tx.Height, &txBody, &tx.TxHash, &tx.TimeStamp, &tx.Code, &tx.Sender)
	if err != nil {
		return types.Tx{}, err
	}

	jsonTx, err := bc.encodeTxBodyJSON(txBody)
	if err != nil {
		return types.Tx{}, err
	}

	tx.Tx = jsonTx

	return tx, nil
}

func (bc tx) GetTxsByHeight(db *sql.DB, height int64, paginate types.Paginate) ([]types.Tx, error) {
	query := `SELECT height,txbody,txhash,timestamp,code,sender 
			  FROM transaction
			  WHERE updated = TRUE AND height = $1 
			  ORDER BY idx DESC
			  OFFSET $2
			  LIMIT $3`
	rows, err := db.Query(query, height, paginate.Offset(), paginate.Limit())
	if err != nil {
		return nil, err
	}

	var txs []types.Tx
	for rows.Next() {
		var tx types.Tx
		var txBody []byte

		err = rows.Scan(&tx.Height, &txBody, &tx.TxHash, &tx.TimeStamp, &tx.Code, &tx.Sender)
		if err != nil {
			return nil, err

		}

		jsonTx, err := bc.encodeTxBodyJSON(txBody)
		if err != nil {
			return nil, err

		}

		tx.Tx = jsonTx
		txs = append(txs, tx)
	}

	return txs, nil
}

func (bc tx) GetTxsByAddress(db *sql.DB, address string, paginate types.Paginate) ([]types.Tx, error) {
	query := `SELECT height,txbody,txhash,timestamp,code,sender 
			  FROM transaction 
			  WHERE updated = TRUE AND sender = $1 
			  ORDER BY idx DESC
			  OFFSET $2
			  LIMIT $3
`
	rows, err := db.Query(query, address, paginate.Offset(), paginate.Limit())
	if err != nil {
		return nil, err
	}

	var txs []types.Tx
	for rows.Next() {
		var tx types.Tx
		var txBody []byte

		err = rows.Scan(&tx.Height, &txBody, &tx.TxHash, &tx.TimeStamp, &tx.Code, &tx.Sender)
		if err != nil {
			return nil, err
		}

		jsonTx, err := bc.encodeTxBodyJSON(txBody)
		if err != nil {
			return nil, err
		}

		tx.Tx = jsonTx
		txs = append(txs, tx)
	}

	return txs, nil
}

func (bc tx) GetTxsCountByAddress(db *sql.DB, address string) (int64, error) {
	query := `SELECT count(*) 
			  FROM transaction 
			  WHERE updated = TRUE AND sender = $1
`
	var count int64
	err := db.QueryRow(query, address).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (bc tx) GetTxsCountByHeight(db *sql.DB, height int64) (int64, error) {
	query := `SELECT count(*) 
			  FROM transaction 
			  WHERE updated = TRUE AND height = $1
`
	var count int64
	err := db.QueryRow(query, height).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (bc tx) encodeTxBodyJSON(txBody []byte) (json.RawMessage, error) {
	txDecoder := klaatoo.MakeEncodingConfig().TxConfig.TxDecoder()
	txJSONEncoder := klaatoo.MakeEncodingConfig().TxConfig.TxJSONEncoder()

	decodedTx, err := txDecoder(txBody)
	if err != nil {
		return nil, err
	}

	encodedJSON, err := txJSONEncoder(decodedTx)
	if err != nil {
		return nil, err
	}

	return encodedJSON, nil
}
