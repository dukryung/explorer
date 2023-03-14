package stream

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/hessegg/nikto-explorer/server/api/event"
	"github.com/hessegg/nikto-explorer/server/api/websocket/conn"
	wstypes "github.com/hessegg/nikto-explorer/server/api/websocket/types"
	wserrors "github.com/hessegg/nikto-explorer/server/api/websocket/types/errors"
	"github.com/hessegg/nikto-explorer/types"
	"github.com/hessegg/nikto-explorer/types/config"
	"github.com/hessegg/nikto-explorer/types/log"
)

type Handler struct {
	ctx   types.Context
	event *event.Handler
	mux   *conn.Multiplexer
	wg    sync.WaitGroup

	config config.ApiConfig
}

func NewStreamHandler(ctx types.Context, config config.ApiConfig, event *event.Handler, mux *conn.Multiplexer) *Handler {
	logger := log.NewLogger("ws/stream", config.Log)

	return &Handler{
		ctx:    ctx.WithLogger(logger),
		event:  event,
		mux:    mux,
		config: config,
		wg:     sync.WaitGroup{},
	}
}

func (s *Handler) duration() time.Duration {
	duration := time.Duration(s.config.Handler.WebSocket.Duration)
	return time.Millisecond * duration
}

func (s *Handler) Run() {
	numCPU := runtime.NumCPU()
	ctx := s.ctx

	ctx.Logger().Info("run stream process", "cpu : ", numCPU)
	for range time.Tick(s.duration()) {

		var eventList []wstypes.WSEvent

		select {
		case <-s.ctx.Context().Done():
			ctx.Logger().Info(s.ctx.Context().Err())
			break
		default:
			eventList = s.event.GetEventList()
			ctx.Logger().Info(fmt.Sprintf("execute events list : %v", len(eventList)))
		}

		executeTime := time.Now()
		worker := 0

		for _, wsEvent := range eventList {
			if worker >= numCPU {
				s.wg.Wait()
				worker = 0
			}

			s.wg.Add(1)
			worker++

			go func(worker int, wsEvent wstypes.WSEvent) {
				ctx := s.ctx
				defer func() {
					if v := recover(); v != nil {
						ctx.Logger().Error("recovered worker : ", worker, v, wsEvent)
						s.wsError(wsEvent, wserrors.Wrap(wserrors.ErrInvalidRequest, wserrors.ErrInternalCode, "internal error"))

						err := s.event.RemoveEvent(wsEvent)
						if err != nil {
							ctx.Logger().Error(err)
						}
					}
					s.wg.Done()
				}()

				bz, err := wsEvent.Execute(ctx)
				if err != nil {
					ctx.Logger().Error(err)
					s.wsError(wsEvent, err)

					err = s.event.RemoveEvent(wsEvent)
					if err != nil {
						ctx.Logger().Error(err)
					}
					return
				}

				s.wsCall(wsEvent, bz)
			}(worker, wsEvent)
		}

		// wait until all worker done
		s.wg.Wait()

		ctx.Logger().Info("execute time : ", time.Now().Sub(executeTime).String())
	}
}

func (s *Handler) wsError(wsEvent wstypes.WSEvent, err error) {
	call := wstypes.WSCall{
		SessionId: wsEvent.SessionId,
		Object:    wstypes.WSError(wsEvent.Request, err),
	}

	s.mux.Call(call)
}

func (s *Handler) wsCall(wsEvent wstypes.WSEvent, v interface{}) {
	call := wstypes.WSCall{
		SessionId: wsEvent.SessionId,
		Object: wstypes.WSResponse{
			Method: wsEvent.Request.Method,
			Code:   0,
			ID:     wsEvent.Request.ID,
			Result: v,
		},
	}

	s.mux.Call(call)
}
