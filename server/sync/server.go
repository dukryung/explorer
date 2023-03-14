package sync

import (
	"time"

	"github.com/hessegg/nikto-explorer/server/sync/service"
	"github.com/hessegg/nikto-explorer/types"
	"github.com/hessegg/nikto-explorer/types/config"
	"github.com/hessegg/nikto-explorer/types/log"
	_ "github.com/lib/pq"
)

var _ types.Server = &Server{}

type Server struct {
	ctx types.Context

	txService        *service.TxService
	blockService     *service.BlockService
	validatorService *service.ValidatorService
	tokenService     *service.TokenService
	nftService       *service.NftService

	serviceManager *types.ServiceManager

	close chan bool
}

func NewServer(config config.SyncConfig) (*Server, error) {
	s := Server{}
	var err error

	logger := log.NewLogger("sync/server", config.Log)
	s.close = make(chan bool)

	db, err := config.DB.GetDBConnection()
	if err != nil {
		logger.Error(err)
	}

	conn, err := config.Node.GetGRPCConnection()
	if err != nil {
		logger.Error(err)
	}

	s.ctx = types.NewContextCancel().
		WithLogger(logger).
		WithDB(db).
		WithGRPCConn(conn)

	s.txService = service.NewTxService(s.ctx, config)
	s.validatorService = service.NewValidatorService(s.ctx, config.Log)
	s.blockService = service.NewBlockService(s.ctx, s.txService, s.validatorService, config)
	s.tokenService = service.NewTokenService(s.ctx, config.Log)
	s.nftService = service.NewNftService(s.ctx, config.Log)

	s.serviceManager = types.NewServiceManager(s.ctx,
		s.txService,
		s.validatorService,
		s.blockService,
		s.tokenService,
		s.nftService,
	)

	return &s, nil
}

func (s *Server) Run() {
	s.ctx.Logger().Info("run server")

	defer func() {
		s.ctx.Close()
	}()

	for range time.Tick(time.Millisecond * 1000) {
		select {
		case <-s.close:
			s.ctx.Logger().Info("shutdown")
			return
		default:
			go s.serviceManager.Run()
		}
	}
}

func (s *Server) Close() {
	s.close <- true
}
