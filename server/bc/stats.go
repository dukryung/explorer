package bc

import (
	"database/sql"
	"time"

	"github.com/hessegg/nikto-explorer/server/bc/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

type stats struct{}

var StatsProvider stats

func (bc stats) GetStats(db *sql.DB) (types.Stats, error) {
	var raw []byte
	query := `SELECT raw FROM block ORDER BY height DESC LIMIT 1`
	err := db.QueryRow(query).Scan(&raw)
	if err != nil {
		return types.Stats{}, err
	}

	latestBlock := &tmproto.Block{}
	err = latestBlock.Unmarshal(raw)
	if err != nil {
		return types.Stats{}, err
	}

	query = `SELECT
				(SELECT COUNT(*) FROM transaction) as tx_total,
				(SELECT COUNT(*) FROM token) as token_total,
				(SELECT COUNT(*) FROM nftoken) as nft_total,
				(SELECT COUNT(*) FROM validator) as validator_total,
				(SELECT	SUM(tokens::numeric)::varchar FROM validator) as total_bonded_tokens,
				(SELECT height from block order by height DESC LIMIT 1) as block_height,
				(SELECT cast(AVG(diff_time) as bigint) FROM block WHERE block_time > $1) as block_avg_time,
				(SELECT cast(AVG(dt) as bigint) FROM (SELECT diff_time as dt FROM block WHERE block_time > $1 ORDER BY diff_time ASC LIMIT 20) as xx) as block_min_time
`

	var stats types.Stats
	err = db.QueryRow(query, latestBlock.Header.Time.Add(-24*time.Hour)).Scan(
		&stats.TxTotal,
		&stats.TokenTotal,
		&stats.NftTotal,
		&stats.ValidatorTotal,
		&stats.TotalBondedTokens,
		&stats.BlockHeight,
		&stats.BlockAvgTime,
		&stats.BlockMinTime,
	)

	if err != nil {
		return types.Stats{}, err
	}

	return stats, nil
}

func (bc stats) GetTxStats(db *sql.DB) ([]types.TxStats, error) {

	query := `SELECT d.date AS time_stamp, COUNT(tx.idx) as tx_count FROM
(SELECT date::date
FROM   GENERATE_SERIES(NOW() - INTERVAL '13' DAY, NOW(), '1 day') date) d
LEFT OUTER JOIN
transaction tx 
ON d.date::date = tx.timestamp::date
GROUP BY d.date
ORDER BY d.date ASC`

	var txStats []types.TxStats
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var txStat types.TxStats
		err := rows.Scan(&txStat.TimeStamp, &txStat.TxCount)
		if err != nil {
			return nil, err
		}

		txStats = append(txStats, txStat)
	}
	return txStats, nil
}
