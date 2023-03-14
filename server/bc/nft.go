package bc

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/hessegg/nikto-explorer/server/bc/types"
	"github.com/lib/pq"
)

type nftProvider struct{}

var NftProvider nftProvider

type WhereClause struct {
	Clause string
	Values []any
}

func (bc nftProvider) GetNftList(db *sql.DB, paginate types.Paginate, where_clauses ...WhereClause) ([]types.NftInfo, error) {
	query_select := `SELECT owner_address, issuer_address, delegated_to_address, block_height, id, url, hash, info_url, info_hash, preview_url, name, collection_id FROM nftoken`

	args := []any{paginate.Offset(), paginate.Limit()}
	clauses := []string{}
	for _, clause := range where_clauses {
		clauses = append(clauses, clause.Clause)
		args = append(args, clause.Values...)
	}

	where := strings.Join(clauses, ") AND (")
	if where != "" {
		where = "\nWHERE (" + where + ")"
	}
	query_pagination := `
				ORDER BY id 
				OFFSET $1
				LIMIT $2`
	query := query_select + where + query_pagination

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	var infos []types.NftInfo
	for rows.Next() {
		var info types.NftInfo

		err = rows.Scan(&info.OwnerAddress, &info.IssuerAddress, &info.DelegatedToAddress, &info.Height,
			&info.Token.Id, &info.Token.Url, &info.Token.Hash, &info.Token.InfoUrl, &info.Token.InfoHash, &info.Token.PreviewUrl, &info.Token.Name, &info.Token.CollectionId)
		if err != nil {
			return nil, err
		}
		infos = append(infos, info)
	}

	return infos, nil
}

func (bc nftProvider) GetNftCount(db *sql.DB, where_clauses ...WhereClause) int64 {
	query_select := `SELECT COUNT(*) as token_count FROM nftoken`

	args := []any{}
	clauses := []string{}
	for _, clause := range where_clauses {
		clauses = append(clauses, clause.Clause)
		args = append(args, clause.Values...)
	}

	where := strings.Join(clauses, ") AND (")
	if where != "" {
		where = "\nWHERE (" + where + ")"
	}
	query_pagination := `
				OFFSET 0
				LIMIT ALL`
	query := query_select + where + query_pagination

	row := db.QueryRow(query, args...)
	if row == nil {
		panic(fmt.Sprintf("%s (%v)", query, args))
		// return 0
	}

	var total int64
	err := row.Scan(&total)
	if err != nil {
		panic(err)
		// return 0
	}

	return total
}

func (bc nftProvider) GetNftCollections(db *sql.DB, paginate types.Paginate) ([]types.NftCollectionInfo, error) {
	subquery := `SELECT COUNT(*) as token_count, collection_id as tcid FROM nftoken
				WHERE collection_id <> ''
				GROUP BY collection_id
				ORDER BY collection_id
				OFFSET $1
				LIMIT $2`
	query := `SELECT C.token_count, owner_address, issuer_address, delegated_to_address, block_height, id, url, hash, info_url, info_hash, preview_url, name, collection_id FROM nftoken
				INNER JOIN (` + subquery + `) as C
				ON (id = C.tcid)`

	rows, err := db.Query(query, paginate.Offset(), paginate.Limit())
	if err != nil {
		return nil, err
	}

	var infos []types.NftCollectionInfo
	for rows.Next() {
		var info types.NftCollectionInfo

		err = rows.Scan(
			&info.Count,
			&info.OwnerAddress, &info.IssuerAddress, &info.DelegatedToAddress, &info.Height,
			&info.Token.Id, &info.Token.Url, &info.Token.Hash, &info.Token.InfoUrl, &info.Token.InfoHash, &info.Token.PreviewUrl, &info.Token.Name, &info.Token.CollectionId)
		if err != nil {
			return nil, err
		}
		infos = append(infos, info)
	}

	return infos, nil
}

func (bc nftProvider) GetNftCollectionCount(db *sql.DB) int64 {
	subquery := `SELECT 1 as dummy FROM nftoken
				WHERE collection_id <> ''
				GROUP BY collection_id
				ORDER BY collection_id`
	query := `SELECT COUNT(*) FROM (` + subquery + `) AS C`

	row := db.QueryRow(query)
	if row == nil {
		return 0
	}

	var total int64
	err := row.Scan(&total)
	if err != nil {
		return 0
	}

	return total
}

func (bc nftProvider) GetNftInfo(db *sql.DB, id string) (types.NftInfo, error) {
	query := `SELECT owner_address, issuer_address, delegated_to_address, block_height, id, url, hash, info_url, info_hash, preview_url, name, collection_id FROM nftoken
			  WHERE id=$1`

	var info types.NftInfo

	err := db.QueryRow(query, id).Scan(&info.OwnerAddress, &info.IssuerAddress, &info.DelegatedToAddress, &info.Height,
		&info.Token.Id, &info.Token.Url, &info.Token.Hash, &info.Token.InfoUrl, &info.Token.InfoHash, &info.Token.PreviewUrl, &info.Token.Name, &info.Token.CollectionId)
	if err != nil {
		return types.NftInfo{}, err
	}

	return info, nil
}

func (bc nftProvider) GetNftInfos(db *sql.DB, ids []string) ([]types.NftInfo, error) {
	query := `SELECT owner_address, issuer_address, delegated_to_address, block_height, id, url, hash, info_url, info_hash, preview_url, name, collection_id FROM nftoken
			  WHERE denom = ANY ($1)`

	rows, err := db.Query(query, pq.Array(ids))
	if err != nil {
		print(err.Error())
		return nil, err
	}

	var infos []types.NftInfo
	for rows.Next() {
		var info types.NftInfo

		err = rows.Scan(&info.OwnerAddress, &info.IssuerAddress, &info.DelegatedToAddress, &info.Height,
			&info.Token.Id, &info.Token.Url, &info.Token.Hash, &info.Token.InfoUrl, &info.Token.InfoHash, &info.Token.PreviewUrl, &info.Token.Name, &info.Token.CollectionId)
		if err != nil {
			return nil, err
		}
		infos = append(infos, info)
	}

	return infos, nil
}
