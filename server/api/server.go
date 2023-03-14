package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hessegg/nikto-explorer/server/api/event"
	"github.com/hessegg/nikto-explorer/server/api/rest"
	"github.com/hessegg/nikto-explorer/server/api/swagger"
	"github.com/hessegg/nikto-explorer/server/api/websocket"
	"github.com/hessegg/nikto-explorer/types"
	"github.com/hessegg/nikto-explorer/types/config"
	"github.com/hessegg/nikto-explorer/types/log"
)

var _ types.Server = &Server{}

type Server struct {
	ctx    types.Context
	router *mux.Router
	srv    *http.Server

	eventHandler     *event.Handler
	webSocketHandler *websocket.Handler
	restHandler      *rest.Handler
	swaggerHandler   *swagger.Handler

	handlerManager *types.HandlerManager

	close chan bool
}

func NewServer(config config.ApiConfig) (*Server, error) {
	s := Server{
		router: mux.NewRouter(),
		close:  make(chan bool),
	}

	logger := log.NewLogger("api/Server", config.Log)

	// make db connection
	db, err := config.DB.GetDBConnection()
	if err != nil {
		logger.Error(err)
	}

	redisClient := config.Redis.GetRedisConnection()

	// make grpc connection
	grpcConn, err := config.Node.GetGRPCConnection()
	if err != nil {
		logger.Error(err)
	}

	s.ctx = types.NewContextCancel().
		WithLogger(logger).
		WithDB(db).
		WithGRPCConn(grpcConn).
		WithRedis(redisClient)

	s.eventHandler = event.NewHandler(s.ctx, config)
	s.webSocketHandler = websocket.NewHandler(s.ctx, s.eventHandler, config)
	s.restHandler = rest.NewHandler(s.ctx, s.eventHandler, config)
	s.swaggerHandler = swagger.NewHandler(s.ctx, config)

	s.handlerManager = types.NewHandlerManager(
		s.eventHandler,
		s.webSocketHandler,
		s.restHandler,
		s.swaggerHandler,
	)

	s.handlerManager.RegisterRoute(s.router)

	s.srv = &http.Server{
		Handler: s.router,
		Addr:    fmt.Sprintf(":%v", config.Port),
	}

	return &s, err
}

func (s *Server) Run() {
	s.ctx.Logger().Info("run server")

	defer func() {
		s.ctx.Close()
	}()

	s.handlerManager.Run()
	s.run()

	<-s.close
	s.ctx.Logger().Info("shutdown")

	if err := s.srv.Shutdown(s.ctx.Context()); err != nil {
		panic(err)
	}
}

func (s *Server) run() {
	go func() {
		err := s.srv.ListenAndServe()
		if err != nil {
			s.ctx.Logger().Error(err)
		}
	}()
}

func (s *Server) Close() {
	s.close <- true
}
