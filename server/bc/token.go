package bc

import (
	"database/sql"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/hessegg/nikto-explorer/server/bc/types"
	"github.com/lib/pq"
)

type token struct{}

var TokenProvider token

func (bc token) GetTokenList(db *sql.DB, paginate types.Paginate) ([]types.Token, error) {
	query := `SELECT owner_address, symbol, description, denom, precision, amount FROM token
				ORDER BY denom 
				OFFSET $1
				LIMIT $2`

	rows, err := db.Query(query, paginate.Offset(), paginate.Limit())
	if err != nil {
		return nil, err
	}

	var tokens []types.Token
	for rows.Next() {
		var token types.Token
		token.Coin = &sdk.Coin{}

		var amount string
		var ok bool
		err = rows.Scan(&token.OwnerAddress, &token.Symbol, &token.Description, &token.Coin.Denom, &token.Precision, &amount)
		token.Coin.Amount, ok = sdk.NewIntFromString(amount)
		if !ok {
			return nil, fmt.Errorf("failed to parse amount")
		}
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}

func (bc token) GetToken(db *sql.DB, denom string) (types.Token, error) {
	query := `SELECT owner_address, symbol, description, denom, precision, amount FROM token
			  WHERE denom=$1`

	var token types.Token
	token.Coin = &sdk.Coin{}

	var amount string
	var ok bool

	err := db.QueryRow(query, denom).Scan(&token.OwnerAddress, &token.Symbol, &token.Description, &token.Coin.Denom, &token.Precision, &amount)
	if err != nil {
		return types.Token{}, err
	}

	token.Coin.Amount, ok = sdk.NewIntFromString(amount)
	if !ok {
		return types.Token{}, fmt.Errorf("failed to parse amount")
	}

	return token, nil
}

func (bc token) GetTokens(db *sql.DB, denoms []string) ([]types.Token, error) {
	query := `SELECT owner_address, symbol, description, denom, precision, amount FROM token
			  WHERE denom = ANY ($1)`

	rows, err := db.Query(query, pq.Array(denoms))
	if err != nil {
		print(err.Error())
		return nil, err
	}

	var tokens []types.Token
	for rows.Next() {
		var token types.Token
		token.Coin = &sdk.Coin{}

		var amount string
		var ok bool
		err = rows.Scan(&token.OwnerAddress, &token.Symbol, &token.Description, &token.Coin.Denom, &token.Precision, &amount)
		token.Coin.Amount, ok = sdk.NewIntFromString(amount)
		if !ok {
			return nil, fmt.Errorf("failed to parse amount")
		}
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}
