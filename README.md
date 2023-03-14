# klaatoo-explorer

A block explorer for Klaatoo network

## Contents

- [Requirements](#requirements)
- [Installation](#Installation)
- [Quick Start](#quick-start)
- [Server](#server)
    - [Sync]()
        - [Fast Sync](#fast-sync)
    - [API]()
        - [WebSocket](#websocket)
        - [REST Api](#rest-api)
        - [Swagger](#swagger)
- [Client](#client)
    - [Serve](#)
    - [Fetch i18n](#fetch-i18n)
    - [Fetch css](#fetch-css)
- [Deploy](#deploy)
    - [From Jenkins](#from-jenkins)

## Requirements

- Requires GO 1.16+
- Requires node.js 14.4+
- Requires Postgresql
- yarn
- docker 20.10+
- docker-compose version 1.29+

## Installation

1. Install Docker Engine
    - Update the apt package index, and install the latest version of Docker Engine

    ```shell
    sudo apt-get update
    sudo apt-get install docker-ce docker-ce-cli containerd.io docker-compose-plugin
    ```

    - To install a specific version of Docker Engine

   ```shell
    apt-cache madison docker-ce
    ```

    - To install a specific version of Docker Engine

   ```shell
    sudo apt-get install docker-ce=<VERSION_STRING> docker-ce-cli=<VERSION_STRING> containerd.io docker-compose-plugin
    ```

    - Verify that Docker Engine is installed
   ```shell
    sudo docker run hello-world
    ```

2. docker-compose
    - Run this command to download the current stable release of Docker Compose

   ```shell
    DOCKER_CONFIG=${DOCKER_CONFIG:-$HOME/.docker}
    mkdir -p $DOCKER_CONFIG/cli-plugins
    curl -SL https://github.com/docker/compose/releases/download/v2.5.0/docker-compose-linux-x86_64 -o $DOCKER_CONFIG/cli-plugins/docker-compose
    ```
    
   - or if you choose to install Compose for all users

   ```shell
   sudo chmod +x /usr/local/lib/docker/cli-plugins/docker-compose
   ```
   - Test the installation
      
   ```shell
    chmod +x $DOCKER_CONFIG/cli-plugins/docker-compose
    ```

   - Install a specific version

   ```shell
   sudo apt-get install docker-compose-plugin=<VERSION_STRING>
   ```

   - Verify that Docker Compose
                      
   ```shell
    sudo docker compose version
   ```

## Quick Start

1. Clone source

```shell
git clone git@github.com:hessegg/klaatoo-explorer.git
cd ./klaatoo-explorer
```

2. Copy config files

- `default_config_server.json` to `config_server.json`
- `default_config_client.json` to `config_client.json`

```shell
cp ./default_config_server.json ./config_server.json
cp ./default_config_client.json ./config_client.json
```

2. Setup your configuration

`default_config_server.json`

```json
{
  "sync": {
    "enable": true,
    "fast_sync": {
      "enable": true,
      "worker": 100
    },
    "log": {
      "enable": true,
      "level": 3
    },
    "node": {
      "api": "node.v2.klaatoo.net:42002",
      "grpc": "node.v2.klaatoo.net:42003"
    },
    "db": {
      "driver_name": "postgres",
      "user": "postgres",
      "password": "qwer1234",
      "db_name": "klaatoo-explorer",
      "db_port": 5432,
      "idle_conn": 2,
      "max_conn": 10
    }
  },
  "api": {
    "enable": true,
    "port": "10001",
    "handler": {
      "event": {
        "cache": true,
        "duration": 1000
      },
      "rest": {
        "enable": true
      },
      "swagger": {
        "enable": true
      },
      "websocket": {
        "enable": true,
        "duration": 1000
      }
    },
    "log": {
      "enable": true,
      "level": 2
    },
    "node": {
      "api": "node.v2.klaatoo.net:42002",
      "grpc": "node.v2.klaatoo.net:42003"
    },
    "db": {
      "driver_name": "postgres",
      "user": "postgres",
      "password": "qwer1234",
      "db_name": "klaatoo-explorer",
      "db_port": 5432,
      "idle_conn": 2,
      "max_conn": 10
    },
    "redis": {
      "host": "localhost",
      "port": 6379,
      "db": 0
    }
  }
}
```

`default_config_client.json`

```json
{
  "network": {
    "localnet": {
      "name": "Local net",
      "api": "http://localhost:8081",
      "ws": "ws://localhost:8081/ws",
      "bech32": {
        "account": "nik",
        "validator": "vl",
        "operator": "op"
      },
      "networkToken": {
        "base": "nnik",
        "symbol": "NIK",
        "precision": 9,
        "description": "Niktonet Network Token"
      }
    },
    "testnet": {
      "name": "Test net",
      "api": "https://explorer.niktonet.com",
      "ws": "wss://explorer.niktonet.com/ws",
      "bech32": {
        "account": "nik",
        "validator": "vl",
        "operator": "op"
      },
      "networkToken": {
        "base": "nnik",
        "symbol": "NIK",
        "precision": 9,
        "description": "Niktonet Network Token"
      }
    },
    "mainnet": {
      "name": "Main net",
      "api": "https://explorer.niktonet.com",
      "ws": "wss://explorer.niktonet.com/ws",
      "bech32": {
        "account": "nik",
        "validator": "vl",
        "operator": "op"
      },
      "networkToken": {
        "base": "nnik",
        "symbol": "NIK",
        "precision": 9,
        "description": "Niktonet Network Token"
      }
    }
  }
}
```

3. Build cli from source

###### Linux

```shell
make build
```

##### Windows

```shell
go build cmd/main.go -o klaatoo-explorer
```

4. Setup postgresql database

```shell
# linux
# cd build

klaatoo-explorer server setup
```

5. Run server

```shell
# linux
# cd build

klaatoo-explorer server run
```

6. Serve client

```shell
# linux
# cd build

klaatoo-explorer client serve
```

## Server

## Fast Sync

1. you can enable fast sync in `default_config_server.json`

> UPDATE fast sync is default. must set enable `true` for sync blocks.

```json
{
  "sync": {
    "enable": true,
    "fast_sync": {
      "enable": true,
      "worker": 8
    }
  }
}
```

- enable
    - true : sync blocks across multiple goroutines.
    - false: sync blocks in single goroutine.
- worker : limit of goroutines.

## WebSocket

```json
{
  "sync": {},
  "api": {
    "handler": {
      "websocket": {
        "enable": true
      }
    },
    "log": {}
  }
}
```

## REST Api

```json
{
  "sync": {},
  "api": {
    "handler": {
      "rest": {
        "enable": true
      }
    },
    "log": {}
  }
}
```

## Swagger

1. install swag

```shell
go get -u github.com/swaggo/swag/cmd/swag

# 1.16 or newer
go install github.com/swaggo/swag/cmd/swag@latest
```

2. you can enable swagger in `default_config_server.json`

```json
{
  "sync": {},
  "api": {
    "handler": {
      "swagger": {
        "enable": true
      }
    },
    "log": {}
  }
}
```

3. build swagger document

##### Linux

```shell
make swagger-gen
```

##### Windows

```shell 
swag init -d ./server --pd --parseDepth 10
```

## Client

## Fetch i18n

1. Open [Link](https://docs.google.com/spreadsheets/d/1aYS_wpVQIq5jvgID6Apdy32964XefyVfbk_JhTiAI0o/edit?usp=sharing) and
   add your `key` and `translation data`

```json
{
  "key": "data"
}
```

2. Fetch i18n data

```shell
/klaatoo-explorer/client$ yarn fetch:i18n
```

3. Apply translation key to `.svelte`

```html

<script>
    import {_} from "svelte-i18n"
</script>

<div>{$_('key')}</div>
```

4. (Optional) Rebuild client

```shell
/klaatoo-explorer/client$ yarn build
```

### Deploy

#### From Jenkins

1. Open [Link](https://jenkins.klaatoo.net/job/klaatoo-explorer/)

2. Click `Build Now`

- [Server Config](https://jenkins.klaatoo.net/configfiles/editConfig?id=config_server)
- [FrontEnd Config](https://jenkins.klaatoo.net/configfiles/editConfig?id=config_client)

#### WebHook

1. Register payload {JENKINS_URL}:{PORT}/github-webhook/
