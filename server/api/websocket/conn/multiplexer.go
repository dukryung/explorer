package conn

import (
	"encoding/json"
	"sync"

	"github.com/hessegg/nikto-explorer/server/api/event"
	wstypes "github.com/hessegg/nikto-explorer/server/api/websocket/types"
	wserrors "github.com/hessegg/nikto-explorer/server/api/websocket/types/errors"
	"github.com/hessegg/nikto-explorer/types"
	"github.com/hessegg/nikto-explorer/types/config"
	"github.com/hessegg/nikto-explorer/types/log"
)

type Multiplexer struct {
	ctx   types.Context
	event *event.Handler

	connMux *sync.Mutex
	conns   map[string]*Conn
	add     chan *Conn
	remove  chan *Conn
	wsCall  chan wstypes.WSCall
}

func NewMultiplexer(ctx types.Context, config config.ApiConfig, event *event.Handler) *Multiplexer {
	return &Multiplexer{
		ctx:   ctx.WithLogger(log.NewLogger("ws/multiplexer", config.Log)),
		event: event,

		add:    make(chan *Conn),
		remove: make(chan *Conn),
		conns:  make(map[string]*Conn),
		wsCall: make(chan wstypes.WSCall),

		connMux: &sync.Mutex{},
	}
}

// Add client connection
func (mux *Multiplexer) Add(conn *Conn) {
	mux.add <- conn
}

// Remove client connection
func (mux *Multiplexer) Remove(conn *Conn) {
	mux.remove <- conn
}

// Call send WSResponse to client by session id
func (mux *Multiplexer) Call(call wstypes.WSCall) {
	//mux.logger.Debug("response ", call)
	mux.wsCall <- call
}

// ConnCount returns current connection count
func (mux *Multiplexer) ConnCount() int {
	return len(mux.conns)
}

func (mux *Multiplexer) Run() {
	for {
		select {
		case <-mux.ctx.Context().Done():
			break
		case conn := <-mux.add:
			mux.addConnection(conn)
		case conn := <-mux.remove:
			mux.removeConnection(conn)
		case v := <-mux.wsCall:
			mux.call(v)
		}
	}
}

func (mux *Multiplexer) call(call wstypes.WSCall) {
	defer mux.connMux.Unlock()
	mux.connMux.Lock()

	if conn, ok := mux.conns[call.SessionId]; ok {
		err := conn.WriteJSON(call.Object)
		if err != nil {
			mux.ctx.Logger().Error(err)
		}
	}
}

func (mux *Multiplexer) removeConnection(conn *Conn) {
	defer mux.connMux.Unlock()
	mux.connMux.Lock()

	if _, ok := mux.conns[conn.Session.Id]; ok {
		mux.ctx.Logger().Info("connection closed :", conn.Session.Id)
		mux.event.RemoveSessionEvents(conn.Session.Id)

		delete(mux.conns, conn.Session.Id)
	}
}

func (mux *Multiplexer) addConnection(conn *Conn) {
	defer mux.connMux.Unlock()
	mux.connMux.Lock()

	mux.conns[conn.Session.Id] = conn
	go mux.HandleConn(conn)
}

func (mux *Multiplexer) broadcastText(message []byte) {
	//mux.logger.Info("broadcast text, conn count : ", mux.ConnCount())
	for _, conn := range mux.conns {
		err := conn.WriteText(message)
		if err != nil {
			mux.removeConnection(conn)
		}
	}
}

func (mux *Multiplexer) broadcastJSON(v interface{}) {
	//mux.logger.Info("broadcast json, conn count : ", mux.ConnCount())
	for _, conn := range mux.conns {
		err := conn.WriteJSON(v)
		if err != nil {
			mux.removeConnection(conn)
		}
	}
}

func (mux *Multiplexer) HandleConn(conn *Conn) {
	defer func() {
		mux.remove <- conn
	}()

	c := conn.Conn

	for {
		_, p, err := c.ReadMessage()
		if err != nil {
			mux.ctx.Logger().Error(err)
			break
		}

		var request wstypes.WSRequest
		err = json.Unmarshal(p, &request)

		// failed to unmarshal request.
		if err != nil {
			mux.ctx.Logger().Error(err)
			err = wserrors.Wrap(wserrors.ErrInvalidRequest, wserrors.ErrInvalidRequestCode, err.Error())

			err = c.WriteJSON(wstypes.WSError(request, err))
			if err != nil {
				mux.ctx.Logger().Error(err)
			}
			continue
		}

		// failed to validate request
		if err = request.Validate(); err != nil {
			mux.ctx.Logger().Error(err)

			err = conn.WriteJSON(wstypes.WSError(request, err))
			if err != nil {
				mux.ctx.Logger().Error(err)
			}
			continue
		}

		// make wsEvent
		wsEvent := wstypes.WSEvent{
			SessionId: conn.Session.Id,
			Request:   request,
		}

		// add to event list
		err = mux.event.Add(wsEvent)

		// failed to add event
		if err != nil {
			err = conn.WriteJSON(wstypes.WSError(wsEvent.Request, err))
			if err != nil {
				mux.ctx.Logger().Error(err)
			}
			continue
		}
	}
}
