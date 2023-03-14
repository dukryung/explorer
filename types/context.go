package types

import (
	"context"
	"database/sql"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/go-redis/redis/v8"
	"github.com/hessegg/nikto-explorer/server/bc/types"
	"github.com/hessegg/nikto-explorer/types/config"
	"github.com/hessegg/nikto-explorer/types/log"
	klaatoo "github.com/hessegg/nikto/app"
	"google.golang.org/grpc"
)

type Context struct {
	ctx         context.Context
	cancel      context.CancelFunc
	db          *sql.DB
	redisClient *redis.Client
	logger      *log.Logger
	paginate    types.Paginate
	grpcClient  *grpc.ClientConn
	txConfig    client.TxConfig
}

func (c Context) Context() context.Context          { return c.ctx }
func (c Context) DB() *sql.DB                       { return c.db }
func (c Context) GRPCConn() *grpc.ClientConn        { return c.grpcClient }
func (c Context) Redis() *redis.Client              { return c.redisClient }
func (c Context) Logger() *log.Logger               { return c.logger }
func (c Context) TxConfig() client.TxConfig         { return c.txConfig }
func (c Context) Value(key interface{}) interface{} { return c.ctx.Value(key) }

func NewContext() Context {
	return Context{
		ctx:      context.Background(),
		logger:   log.NewLogger("ctx", config.DefaultLogConfig()),
		txConfig: klaatoo.MakeEncodingConfig().TxConfig,
		cancel:   func() {},
	}
}

func NewContextCancel() Context {
	ctx, cancel := context.WithCancel(context.Background())
	return Context{
		ctx:      ctx,
		logger:   log.NewLogger("ctx", config.DefaultLogConfig()),
		txConfig: klaatoo.MakeEncodingConfig().TxConfig,
		cancel:   cancel,
	}
}

func (c Context) WithContext(ctx context.Context) Context {
	c.ctx = ctx
	return c
}

func (c Context) WithCancel(cancel context.CancelFunc) Context {
	c.cancel = cancel
	return c
}

func (c Context) WithDB(db *sql.DB) Context {
	if db == nil {
		c.Logger().Error("db is nil")
	}
	c.db = db
	return c
}

func (c Context) WithGRPCConn(conn *grpc.ClientConn) Context {
	if conn == nil {
		c.Logger().Error("grpc conn is nil")
	}
	c.grpcClient = conn
	return c
}

func (c Context) WithRedis(conn *redis.Client) Context {
	if conn == nil {
		c.Logger().Error("redis conn is nil")
	}
	c.redisClient = conn
	return c
}

func (c Context) WithLogger(logger *log.Logger) Context {
	c.logger = logger
	return c
}

func (c Context) WithValue(key, value interface{}) Context {
	c.ctx = context.WithValue(c.ctx, key, value)
	return c
}

func (c Context) Close() {
	c.cancel()

	if c.db != nil {
		err := c.db.Close()
		if err != nil && c.logger != nil {
			c.logger.Error(err)
		}
	}
	if c.grpcClient != nil {
		err := c.grpcClient.Close()
		if err != nil && c.logger != nil {
			c.logger.Error(err)
		}
	}
}

const ContextKey = "context"

func WrapContext(ctx Context) context.Context {
	return context.WithValue(ctx.ctx, ContextKey, ctx)
}

func UnwrapContext(ctx context.Context) Context {
	return ctx.Value(ContextKey).(Context)
}
