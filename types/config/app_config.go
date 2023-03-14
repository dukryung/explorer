package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"os"
	"runtime"
)

const (
	DefaultServerConfigPath = "./default_config_server.json"
	ServerConfigPath        = "./config_server.json"

	DefaultClientConfigPath = "./default_config_client.json"
	ClientConfigPath        = "./config_client.json"
)

var logger = log.New(os.Stdout, color.CyanString("INFO "), log.Ltime|log.Lshortfile)

type AppConfig struct {
	Sync SyncConfig `json:"sync"`
	Api  ApiConfig  `json:"api"`
}

type SyncConfig struct {
	DefaultConfig

	FastSync FastSyncConfig `json:"fast_sync"`
	Node     NodeConfig     `json:"node"`
	Log      LogConfig      `json:"log"`
	DB       DBConfig       `json:"db"`
}

type ApiConfig struct {
	DefaultConfig

	Port    string        `json:"port"`
	Handler HandlerConfig `json:"handler"`
	Node    NodeConfig    `json:"node"`
	Log     LogConfig     `json:"log"`
	DB      DBConfig      `json:"db"`
	Redis   RedisConfig   `json:"redis"`
}

//--------------- Sync Configs

type FastSyncConfig struct {
	DefaultConfig
	Worker int64 `json:"worker"`
}

//--------------- Handler Configs

type HandlerConfig struct {
	Event     EventConfig     `json:"event"`
	WebSocket WebSocketConfig `json:"websocket"`
	REST      RESTConfig      `json:"rest"`
	Swagger   SwaggerConfig   `json:"swagger"`
}

type EventConfig struct {
	Cache    bool  `json:"cache"`
	Duration int64 `json:"duration"`
}

type WebSocketConfig struct {
	DefaultConfig
	Worker   int64 `json:"worker"`
	Duration int64 `json:"duration"`
}

type SwaggerConfig struct {
	DefaultConfig
}

type RESTConfig struct {
	DefaultConfig
}

//-------------------- Common Configs

type DefaultConfig struct {
	Enable bool `json:"enable"`
}

type LogConfig struct {
	Enable bool `json:"enable"`
	Level  int  `json:"level"`
}

type DBConfig struct {
	Host       string `json:"host"`
	DriverName string `json:"driver_name"`
	User       string `json:"user"`
	Password   string `json:"password"`
	DBName     string `json:"db_name"`
	DBPort     int    `json:"db_port"`
	IdleConn   int    `json:"idle_conn"`
	MaxConn    int    `json:"max_conn"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Password string `json:"password"`
	Port     int    `json:"port"`
}

type NodeConfig struct {
	API  string `json:"api"`
	GRPC string `json:"grpc"`
}

func (config NodeConfig) GetGRPCConnection() (*grpc.ClientConn, error) {
	return grpc.Dial(config.GRPC, grpc.WithInsecure())
}

func (config DBConfig) GetDBInfo() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", config.Host, config.User, config.Password, config.DBName, config.DBPort)
}

func (config DBConfig) GetDBConnection() (*sql.DB, error) {
	db, err := sql.Open(config.DriverName, config.GetDBInfo())
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(config.IdleConn)
	db.SetMaxOpenConns(config.MaxConn)
	return db, nil
}

func (config RedisConfig) GetRedisConnection() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       0,
	})

	return redisClient
}

func (config *AppConfig) UnmarshalJSON(data []byte) error {
	type tempAppConfig AppConfig
	appConfig := tempAppConfig{
		Sync: SyncConfig{
			FastSync: FastSyncConfig{
				Worker: int64(runtime.NumCPU()),
			},
			Log: DefaultLogConfig(),
			DB: DBConfig{
				Host: "localhost",
			},
		},
		Api: ApiConfig{
			Handler: HandlerConfig{
				Event: EventConfig{
					Cache:    true,
					Duration: 1000,
				},
				WebSocket: WebSocketConfig{
					Duration: 1000,
					Worker:   int64(runtime.NumCPU()),
				},
			},
			Log: DefaultLogConfig(),
			DB: DBConfig{
				Host: "localhost",
			},
		},
	}

	err := json.Unmarshal(data, &appConfig)
	if err != nil {
		return err
	}

	*config = AppConfig(appConfig)
	return nil
}

func (config *AppConfig) LoadConfig(path string) error {
	if path == "" {
		path = ServerConfigPath
	}

	logger.Output(3, fmt.Sprint("load config", path))

	data, err := ioutil.ReadFile(path)
	if err != nil {
		data, err = ioutil.ReadFile(DefaultServerConfigPath)
		if err != nil {
			panic(fmt.Sprintf("failed to read config :\t\n%v", path))
		}
	}

	err = config.UnmarshalJSON(data)
	if err != nil {
		panic(fmt.Sprintf("can not unmarshal config :\t\n %v", err))
	}

	logger.Output(3, fmt.Sprint("loaded config", config.String()))

	return nil
}

func DefaultLogConfig() LogConfig {
	return LogConfig{
		Enable: true,
		Level:  1,
	}
}

func (config *AppConfig) String() string {
	bz, _ := json.MarshalIndent(config, "", "\t")
	return string(bz)
}
