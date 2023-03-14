package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
	"github.com/hessegg/nikto-explorer/server/api/event"
	"github.com/hessegg/nikto-explorer/server/api/event/codec"
	"github.com/hessegg/nikto-explorer/types"
	"github.com/hessegg/nikto-explorer/types/config"
)

type Handler struct {
	ctx    types.Context
	config config.ApiConfig
	event  *event.Handler
}

var _ types.Handler = &Handler{}

func NewHandler(ctx types.Context, event *event.Handler, config config.ApiConfig) *Handler {
	return &Handler{
		ctx:    ctx,
		config: config,
		event:  event,
	}
}

func (h *Handler) handleFunc(w http.ResponseWriter, r *http.Request) {
	types.AllowCORS(w)

	path := r.URL.Path
	funcs := h.event.GWEventMethods()
	params := codec.Decode(path, r.URL.Query())
	h.ctx.Logger().Info(path, params)

	if f, ok := funcs[path]; ok {
		ctx := types.NewContext().
			WithContext(h.ctx.Context()).
			WithDB(h.ctx.DB()).
			WithGRPCConn(h.ctx.GRPCConn()).
			WithRedis(h.ctx.Redis())

		res, err := h.call(*f, ctx, params...)
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			// the first response is byte array
			bz := res[0].Interface()
			// the second response is error
			if !res[1].IsNil() {
				switch res[1].Interface().(type) {
				case *error:

				default:
					w.Write([]byte("{message:\"internal error\"}"))
					return
				}
			}

			json, err := json.Marshal(bz)
			if err != nil {
				w.Write([]byte(err.Error()))
			} else {
				w.Write(json)
			}
		}
	} else {
		w.Write([]byte("invalid endpoint"))
	}
}

func (h *Handler) call(function interface{}, ctx types.Context, params ...interface{}) ([]reflect.Value, error) {
	f := reflect.ValueOf(function)

	in := []reflect.Value{
		reflect.ValueOf(ctx),
	}
	for _, param := range params {
		in = append(in, reflect.ValueOf(param))
	}

	if len(in) != f.Type().NumIn() {
		return nil, fmt.Errorf("invalid params")
	}

	return f.Call(in), nil
}

func (h *Handler) Run() {}

func (h *Handler) restConfig() config.RESTConfig {
	return h.config.Handler.REST
}

func (h *Handler) RegisterRoute(router *mux.Router) *mux.Router {
	if h.restConfig().Enable {
		for k, _ := range h.event.GWEventMethods() {
			h.ctx.Logger().Debug("RegisterRoute", k)

			router.HandleFunc(k, h.handleFunc)
		}
	}

	return router
}
