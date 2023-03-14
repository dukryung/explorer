package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/hessegg/nikto-explorer/server/app"
	"github.com/hessegg/nikto-explorer/types/config"
)

// @title           Niktonet Explorer - RESTConfig Api
// @version         1.0
// @description     RESTConfig api spec

// @host      https://explorer.niktonet.com/
// @BasePath  /

// @securityDefinitions.basic  BasicAuth
func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	app := app.NewApp(config.ServerConfigPath)
	err := app.RunServers()
	if err != nil {
		panic(err)
	}

	<-quit

	app.CloseServers()
}
