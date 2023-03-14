package types

import (
	"github.com/cosmos/cosmos-sdk/types/query"
	"math"
)

type Paginate struct {
	limit int `json:"limit"`
	page  int `json:"page"`
}

func NewPaginate(limit, page float64) Paginate {
	paginate := Paginate{
		limit: int(math.Round(limit)),
		page:  int(math.Round(page)),
	}

	if paginate.limit > MaxResponseCount {
		paginate.limit = MaxResponseCount
	}

	return paginate
}

func DefaultPaginate() Paginate {
	return Paginate{
		limit: MaxResponseCount,
		page:  0,
	}
}

func (p Paginate) Offset() int {
	return p.limit * p.page
}

func (p Paginate) Limit() int {
	return p.limit
}

func (p Paginate) PageRequest() *query.PageRequest {
	return &query.PageRequest{
		Limit:      uint64(p.limit),
		Offset:     uint64(p.page),
		CountTotal: true,
	}
}
