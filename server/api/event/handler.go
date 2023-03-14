package event

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/hessegg/nikto-explorer/server/api/event/controller"
	wstypes "github.com/hessegg/nikto-explorer/server/api/websocket/types"
	"github.com/hessegg/nikto-explorer/server/api/websocket/types/errors"
	"github.com/hessegg/nikto-explorer/types"
	"github.com/hessegg/nikto-explorer/types/config"
	"github.com/hessegg/nikto-explorer/types/log"
)

const (
	wsEventPath = "%s.%s"
	gwEventPath = "/api/%s/%s"
)

type Handler struct {
	ctx types.Context

	block     controller.Block
	tx        controller.Tx
	token     controller.Token
	stats     controller.Stats
	account   controller.Account
	validator controller.Validator
	nft       controller.Nft

	config         config.ApiConfig
	wsEventMux     sync.Mutex
	wsEvents       map[string]wstypes.WSEvent
	wsEventMethods map[string]*interface{}

	gwEventMethods map[string]*interface{}
}

var _ types.Handler = &Handler{}

func NewHandler(ctx types.Context, config config.ApiConfig) *Handler {
	handler := Handler{}

	handler.block = controller.Block{}

	handler.wsEventMux = sync.Mutex{}
	handler.config = config
	handler.wsEvents = make(map[string]wstypes.WSEvent)
	handler.wsEventMethods = make(map[string]*interface{})
	handler.gwEventMethods = make(map[string]*interface{})
	handler.ctx = ctx.WithLogger(log.NewLogger("ws/event", config.Log))

	handler.init()

	return &handler
}

func (s *Handler) duration() time.Duration {
	duration := time.Duration(s.config.Handler.Event.Duration)
	return time.Millisecond * duration
}

func (s *Handler) Run() {
	for range time.Tick(s.duration()) {
		select {
		case <-s.ctx.Context().Done():
			s.ctx.Logger().Info(s.ctx.Context().Err())
			break
		default:
			if err := s.setCache(); err != nil {
				s.ctx.Logger().Error(err)
			}
		}
	}
}

func (s *Handler) setCache() error {
	s.ctx.Logger().Debug("set cache data")

	err := controller.Cache.SetTxCache(s.ctx)
	if err != nil {
		return err
	}

	err = controller.Cache.SetBlockCache(s.ctx)
	if err != nil {
		return err
	}

	err = controller.Cache.SetStatsCache(s.ctx)
	if err != nil {
		return err
	}

	err = controller.Cache.SetTxStatsCache(s.ctx)
	if err != nil {
		return err
	}

	err = controller.Cache.SetTokenCache(s.ctx)
	if err != nil {
		return err
	}

	err = controller.Cache.SetNftCache(s.ctx)
	if err != nil {
		return err
	}

	err = controller.Cache.SetValidatorCache(s.ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Handler) init() {
	element := reflect.ValueOf(s).Elem()

	for i := 0; i < element.NumField(); i++ {
		field := element.Field(i)
		fieldType := field.Type()
		if fieldType == reflect.TypeOf(sync.Mutex{}) || fieldType == reflect.TypeOf(types.Context{}) {
			continue
		}

		fieldName := strings.ToLower(fieldType.Name())

		value := reflect.New(fieldType)

		for j := 0; j < value.NumMethod(); j++ {
			method := value.Type().Method(j)
			methodName := strings.ToLower(method.Name)

			wsMethodPath := fmt.Sprintf(wsEventPath, fieldName, methodName)
			gwMethodPath := fmt.Sprintf(gwEventPath, fieldName, methodName)

			methodFunc := value.MethodByName(method.Name).Interface()
			s.ctx.Logger().Debug("eventHandler.init: wsMethodPath", wsMethodPath)
			s.wsEventMethods[wsMethodPath] = &methodFunc
			s.ctx.Logger().Debug("eventHandler.init: gwMethodPath", gwMethodPath)
			s.gwEventMethods[gwMethodPath] = &methodFunc
		}
	}
}

func (s *Handler) WSEventMethods() map[string]*interface{} { return s.wsEventMethods }
func (s *Handler) GWEventMethods() map[string]*interface{} { return s.gwEventMethods }

func (s *Handler) Add(wsEvent wstypes.WSEvent) error {
	s.wsEventMux.Lock()
	defer s.wsEventMux.Unlock()

	switch wsEvent.Request.Type {
	case wstypes.Unsubscribe:
		err := s.removeEvent(wsEvent)
		if err != nil {
			return errors.Wrap(err, errors.ErrInternalCode, "internal errors")
		}
		return nil
	default:
		f, ok := s.wsEventMethods[wsEvent.Request.Method]
		if !ok {
			return errors.Wrap(errors.ErrInvalidMethod, errors.ErrInvalidMethodCode, wsEvent.Request.Method)
		}

		wsEvent.Func = f

		//bz, err := wsEvent.Execute(s.ctx)
		//if err != nil {
		//	return errors.Wrap(err, errors.ErrInternalCode, "internal errors")
		//}

		s.ctx.Logger().Info("add event", wsEvent)
		s.wsEvents[wsEvent.EventKey()] = wsEvent
	}

	return nil
}

func (s *Handler) removeEvent(wsEvent wstypes.WSEvent) error {
	s.ctx.Logger().Debug("remove event", wsEvent)

	if _, ok := s.wsEvents[wsEvent.EventKey()]; !ok {
		return fmt.Errorf(fmt.Sprintf("not exist event: %v", wsEvent.EventKey()))
	}
	delete(s.wsEvents, wsEvent.EventKey())

	return nil
}

func (s *Handler) RemoveEvent(wsEvent wstypes.WSEvent) error {
	s.wsEventMux.Lock()
	defer s.wsEventMux.Unlock()

	if err := s.removeEvent(wsEvent); err != nil {
		return err
	}

	return nil
}

func (s *Handler) RemoveSessionEvents(sessionId string) {
	s.wsEventMux.Lock()
	defer s.wsEventMux.Unlock()

	for _, wsEvent := range s.wsEvents {
		if wsEvent.SessionId == sessionId {
			err := s.removeEvent(wsEvent)
			if err != nil {
				s.ctx.Logger().Error(err)
				return
			}
		}
	}
}

func (s *Handler) GetEventList() []wstypes.WSEvent {
	s.wsEventMux.Lock()
	defer s.wsEventMux.Unlock()

	var eventList []wstypes.WSEvent

	for _, wsEvent := range s.wsEvents {
		eventList = append(eventList, wsEvent)

		switch wsEvent.Request.Type {
		case wstypes.SingleRequest:
			err := s.removeEvent(wsEvent)
			if err != nil {
				s.ctx.Logger().Error("internal error", err.Error())
			}
		}
	}

	return eventList
}

func (s *Handler) RegisterRoute(router *mux.Router) *mux.Router {
	return router
}
