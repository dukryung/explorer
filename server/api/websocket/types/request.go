package types

import (
	"fmt"

	bctypes "github.com/hessegg/nikto-explorer/server/bc/types"
)

const (
	SingleRequest = "SINGLE_REQUEST"
	Subscribe     = "SUBSCRIBE"
	Unsubscribe   = "UNSUBSCRIBE"
	Error         = "ERROR"
)

const (
	WSEventKey = "%s/%s/%d"
)

type WSRequest struct {
	Type     string           `json:"type"`
	ID       int              `json:"id"`
	Method   string           `json:"method"`
	Paginate bctypes.Paginate `json:"paginate"`
	Params   []interface{}    `json:"params"`
}

func (r WSRequest) Validate() error {
	if r.Type == "" {
		return fmt.Errorf("empty request type")
	}

	if r.Type != SingleRequest && r.Type != Subscribe && r.Type != Unsubscribe {
		return fmt.Errorf("invalid request type %v", r.Type)
	}

	if r.Method == "" {
		return fmt.Errorf("empty request method")
	}

	if r.ID == 0 {
		return fmt.Errorf("invalid request id")
	}

	return nil
}
