package types

import (
	"fmt"
	"reflect"

	wserrors "github.com/hessegg/nikto-explorer/server/api/websocket/types/errors"
	"github.com/hessegg/nikto-explorer/types"
)

type WSEvent struct {
	SessionId string
	Request   WSRequest
	Func      *interface{}
}

func (e WSEvent) Execute(ctx types.Context) (interface{}, error) {
	req := []reflect.Value{
		reflect.ValueOf(ctx),
	}

	for _, param := range e.Request.Params {
		req = append(req, reflect.ValueOf(param))
	}

	ctx.Logger().Debug("execute event", e.Request.Method)
	// event call returns wrapped error
	res, err := e.call(*e.Func, req)
	if err != nil {
		return nil, err
	}

	// the first response is byte array
	bz := res[0].Interface()
	// the second response is error
	if !res[1].IsNil() {
		// return wrapped error
		return nil, wserrors.Wrap(res[1].Interface().(error), wserrors.ErrInvalidResponseCode, "internal errors")
	}

	return bz, nil
}

func (e WSEvent) EventKey() string {
	return fmt.Sprintf(WSEventKey, e.SessionId, e.Request.Method, e.Request.ID)
}

func (e WSEvent) call(function interface{}, params []reflect.Value) ([]reflect.Value, error) {
	f := reflect.ValueOf(function)

	if len(params) != f.Type().NumIn() {
		err := wserrors.Wrap(wserrors.ErrInvalidParams, wserrors.ErrInvalidParamsCode, "the number of params is not adapted")
		return nil, err
	}

	result := f.Call(params)
	return result, nil
}
