package types

import (
	wserrors "github.com/hessegg/nikto-explorer/server/api/websocket/types/errors"
	"github.com/pkg/errors"
)

func WSError(wsRequest WSRequest, err error) WSResponse {
	switch err := errors.Cause(err).(type) {
	case *wserrors.WrappedError:
		return WSResponse{
			Method: wsRequest.Method,
			ID:     wsRequest.ID,
			Code:   err.Code(),
			Result: err.Error(),
		}
	default:
		return WSResponse{
			Method: wsRequest.Method,
			ID:     wsRequest.ID,
			Code:   wserrors.ErrUnknownCode,
			Result: err.Error(),
		}
	}
}
