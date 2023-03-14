package main

import (
	"os"

	"github.com/hessegg/nikto-explorer/cmd/command"
)

func main() {
	if err := command.NewRootCmd().Execute(); err != nil {
		panic(err)
		os.Exit(1)
	}
}
