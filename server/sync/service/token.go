package service

import (
	"sync"

	"github.com/hessegg/nikto-explorer/types"
	"github.com/hessegg/nikto-explorer/types/config"
	"github.com/hessegg/nikto-explorer/types/log"
	bankz "github.com/hessegg/nikto/x/bankz/types"
)

type TokenService struct {
	ctx   types.Context
	mutex sync.Mutex

	bankzClient bankz.QueryClient
}

var _ types.Service = &TokenService{}

func NewTokenService(ctx types.Context, config config.LogConfig) *TokenService {
	return &TokenService{
		ctx:         ctx.WithLogger(log.NewLogger("sync/token", config)),
		mutex:       sync.Mutex{},
		bankzClient: bankz.NewQueryClient(ctx.GRPCConn()),
	}
}

func (s *TokenService) Run() error {
	if err := s.syncTokenList(); err != nil {
		return err
	}
	return nil
}

func (s *TokenService) syncTokenList() error {
	s.ctx.Logger().Debug("sync token list...")

	res, err := s.bankzClient.List(s.ctx.Context(), &bankz.QueryListRequest{})
	if err != nil {
		return err
	}

	db, err := s.ctx.DB().Begin()
	if err != nil {
		return err
	}

	for _, token := range res.Tokens {
		denomResponse, err := s.bankzClient.NormalizedDenom(
			s.ctx.Context(),
			&bankz.QueryNormalizedDenomRequest{Symbol: token.Symbol},
		)
		if err != nil {
			s.ctx.Logger().Error("failed to query denominator:", err.Error())
			continue
		}
		if denomResponse == nil {
			s.ctx.Logger().Error("empty denominator response")
			continue
		}
		query := `INSERT INTO token (owner_address, symbol, description, precision, denom, amount) 
					VALUES($1, $2, $3, $4, $5, $6)
					ON CONFLICT (denom)
	   				DO UPDATE 
	   				SET owner_address = $1,
	   	   				symbol = $2,
		   				description = $3,
		   				precision = $4,
		   				denom = $5,
		   				amount = $6`

		_, err = db.Exec(query, token.OwnerAddress, token.Symbol, token.Description, token.Precision, denomResponse.Denom, "0")
		if err != nil {
			s.ctx.Logger().Error("INSERT failed: ", err.Error())
			if err := db.Rollback(); err != nil {
				s.ctx.Logger().Error("rollback failed: ", err.Error())
				return err
			}
			return err
		}
	}
	if err := db.Commit(); err != nil {
		s.ctx.Logger().Error("commit failed: ", err.Error())
		return err
	}

	return nil
}
