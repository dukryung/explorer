#!/usr/bin/make -f

BUILDDIR ?= $(CURDIR)/build
CLIENTCONFIG ?= $(CURDIR)/config_client.json
SERVERCONFIG ?= $(CURDIR)/config_server.json

clean:
	rm -rf $(BUILDDIR)/*

build:	clean
	go build -o $(BUILDDIR)/nikto-explorer -ldflags="-s -w" ./cmd/main.go
	ln -s $(CURDIR)/client $(BUILDDIR)/client
	cp $(CLIENTCONFIG) $(BUILDDIR)
	cp $(SERVERCONFIG) $(BUILDDIR)

swagger-gen:
	swag init -d ./server --pd --parseDepth 10
