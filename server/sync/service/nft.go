package service

import (
	"encoding/base64"
	"fmt"
	"sync"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/hessegg/nikto-explorer/types"
	"github.com/hessegg/nikto-explorer/types/config"
	"github.com/hessegg/nikto-explorer/types/log"
	nft "github.com/hessegg/nikto/x/nft/types"
)

type NftService struct {
	ctx   types.Context
	mutex sync.Mutex

	nftClient nft.QueryClient
}

var _ types.Service = &NftService{}

func NewNftService(ctx types.Context, config config.LogConfig) *NftService {
	return &NftService{
		ctx:       ctx.WithLogger(log.NewLogger("sync/nft", config)),
		mutex:     sync.Mutex{},
		nftClient: nft.NewQueryClient(ctx.GRPCConn()),
	}
}

func (s *NftService) Run() error {
	if err := s.syncNftList(); err != nil {
		return err
	}
	return nil
}

func (s *NftService) getMaxHeight() (height int64, empty bool) {
	query := `SELECT COUNT(*) as count, coalesce(MAX(block_height), 0) as height FROM nftoken_base`
	db := s.ctx.DB()
	row := db.QueryRow(query)
	if row == nil {
		panic("Query failed")
	}
	var cnt int64
	err := row.Scan(&cnt, &height)
	if err != nil {
		panic(err)
	}
	return height, (cnt == 0)
}

func (s *NftService) syncNftList() error {
	s.ctx.Logger().Debug("sync nft list...")

	height, empty := s.getMaxHeight()
	if empty {
		height = 0
	} else {
		// start after the max height
		height += 1
	}

	// Sync tokens
	s.ctx.Logger().Debug("sync nft list: start height: ", height)
	res, err := s.nftClient.NFTSince(s.ctx.Context(), &nft.QueryNFTSinceRequest{
		StartHeight: height,
		Pagination: &query.PageRequest{
			CountTotal: true,
			Limit:      1,
		},
	})

	if err != nil {
		s.ctx.Logger().Error("NFTSince query failed: ", err.Error())
		return err
	}
	if res.Pagination == nil {
		s.ctx.Logger().Error("NFTSince query: no pagination")
		return fmt.Errorf("NFTSince query: The server returned no pagination information")
	}
	total := res.Pagination.Total
	s.ctx.Logger().Debug("sync nft list: number of tokens: ", total)

	if total == 0 {
		return nil
	}

	start := uint64(0)
	limit := uint64(100)
	db, err := s.ctx.DB().Begin()
	if err != nil {
		return err
	}

	for ; start < total; start += limit {
		res, err := s.nftClient.NFTSince(s.ctx.Context(), &nft.QueryNFTSinceRequest{
			StartHeight: height,
			Pagination: &query.PageRequest{
				CountTotal: false,
				Offset:     start,
				Limit:      limit,
			},
		})
		if err != nil {
			db.Rollback()
			s.ctx.Logger().Error("NFTSince query failed: ", err.Error())
			return err
		}
		s.ctx.Logger().Debug("sync nft list: processing tokens: ", start, start+limit)

		for _, nfToken := range res.TokenInfos {
			s.ctx.Logger().Debug("sync nft list: processing token: ", nfToken)
			query := `INSERT INTO nftoken_base (owner_address, issuer_address, delegated_to_address, id, url, hash, info_url, info_hash, preview_url, name, collection_id, block_height, burnt) 
					VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
					ON CONFLICT (id)
	   				DO UPDATE 
	   				SET owner_address= $1,
					   issuer_address= $2,
					   delegated_to_address= $3,
					   id= $4,
					   url= $5,
					   hash= $6,
					   info_url= $7,
					   info_hash= $8,
					   preview_url= $9,
					   name = $10,
					   collection_id = $11,
					   block_height = $12,
					   burnt = $13
					   `

			_, err = db.Exec(query,
				nfToken.OwnerAddress, nfToken.IssuerAddress, nfToken.DelegatedToAddress,
				nfToken.Token.Id,
				nfToken.Token.Url, base64.StdEncoding.EncodeToString(nfToken.Token.Hash),
				nfToken.Token.InfoUrl, base64.StdEncoding.EncodeToString(nfToken.Token.InfoHash),
				nfToken.Token.PreviewUrl,
				nfToken.Token.Name,
				nfToken.Token.CollectionId,
				nfToken.Height,
				nfToken.OwnerAddress == "",
			)
			if err != nil {
				s.ctx.Logger().Error("INSERT failed: ", err.Error())
				if err := db.Rollback(); err != nil {
					s.ctx.Logger().Error("rollback failed: ", err.Error())
					return err
				}
				return err
			}
		}
	}

	if err := db.Commit(); err != nil {
		s.ctx.Logger().Error("commit failed: ", err.Error())
		return err
	}

	return nil
}
