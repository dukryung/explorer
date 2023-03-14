package websocket

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/hessegg/nikto-explorer/server/api/event"
	"github.com/hessegg/nikto-explorer/server/api/websocket/conn"
	"github.com/hessegg/nikto-explorer/server/api/websocket/stream"
	"github.com/hessegg/nikto-explorer/types"
	"github.com/hessegg/nikto-explorer/types/config"

	"net/http"
)

const (
	BufferSize = 1024
)

type Handler struct {
	ctx           types.Context
	config        config.ApiConfig
	upgrader      websocket.Upgrader
	connMux       *conn.Multiplexer
	streamHandler *stream.Handler
}

var _ types.Handler = &Handler{}

func NewHandler(ctx types.Context, event *event.Handler, config config.ApiConfig) *Handler {
	handler := Handler{
		ctx:    ctx,
		config: config,
	}

	// initialize multiplexer
	handler.connMux = conn.NewMultiplexer(ctx, config, event)

	// initialize stream handler
	handler.streamHandler = stream.NewStreamHandler(ctx, config, event, handler.connMux)

	handler.upgrader = websocket.Upgrader{
		ReadBufferSize:  BufferSize,
		WriteBufferSize: BufferSize,
	}
	handler.upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	return &handler
}

func (h *Handler) Run() {
	go h.connMux.Run()
	go h.streamHandler.Run()
}

func (h *Handler) webSocketConfig() config.WebSocketConfig {
	return h.config.Handler.WebSocket
}

func (h *Handler) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	types.AllowCORS(w)

	c, err := h.upgrader.Upgrade(w, r, nil)

	if err != nil {
		return
	}

	connection := conn.NewConn(c, w, r)
	h.ctx.Logger().Info("new connection :", connection.Session.Id)
	h.connMux.Add(connection)
}

func (h *Handler) RegisterRoute(router *mux.Router) *mux.Router {
	if h.webSocketConfig().Enable {
		router.HandleFunc("/ws", h.WebSocketHandler)
		h.ctx.Logger().Debug("RegisterRoute", "/ws")
	}

	return router
}
