package types

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Handler interface {
	Run()
	RegisterRoute(router *mux.Router) *mux.Router
}

type HandlerManager struct {
	handlers []Handler
}

func AllowCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func NewHandlerManager(handlers ...Handler) *HandlerManager {
	return &HandlerManager{
		handlers: handlers,
	}
}

func (m *HandlerManager) Run() {
	for _, handler := range m.handlers {
		go handler.Run()
	}
}

func (m *HandlerManager) RegisterRoute(router *mux.Router) *mux.Router {
	for _, handler := range m.handlers {
		handler.RegisterRoute(router)
	}
	return router
}