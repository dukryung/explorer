package app

import (
	"fmt"

	"github.com/hessegg/nikto-explorer/server/api"
	"github.com/hessegg/nikto-explorer/server/sync"
	"github.com/hessegg/nikto-explorer/types"
	"github.com/hessegg/nikto-explorer/types/config"
)

type App struct {
	servers []types.Server
	config  config.AppConfig
}

func NewApp(configPath string) *App {
	app := App{}

	config.SealConfig()
	app.config = config.AppConfig{}
	err := app.config.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}

	app.initServers()

	return &app
}

func (app *App) initServers() {
	if app.config.Sync.Enable {
		syncServer, err := sync.NewServer(app.config.Sync)
		if err != nil {
			panic(err)
		}

		app.servers = append(app.servers, syncServer)
	}

	if app.config.Api.Enable {
		wsServer, err := api.NewServer(app.config.Api)
		if err != nil {
			panic(err)
		}
		app.servers = append(app.servers, wsServer)
	}

}

func (app *App) RunServers() error {
	if len(app.servers) == 0 {
		return fmt.Errorf("no server enabled")
	}

	for _, server := range app.servers {
		go server.Run()
	}

	return nil
}

func (app *App) CloseServers() {
	for _, server := range app.servers {
		server.Close()
	}
}
