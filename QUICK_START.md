# Block Explorer for EXAIS blockchain

A block explorer for Klaatoo network

## Contents

- [Block Explorer for EXAIS blockchain](#block-explorer-for-exais-blockchain)
  - [Contents](#contents)
  - [Requirements](#requirements)
  - [Installation](#installation)
  - [Quick Start](#quick-start)
    - [Unpack the archive](#unpack-the-archive)
    - [Review the config files](#review-the-config-files)
      - [Backend Configuration](#backend-configuration)
      - [Frontend Configuration](#frontend-configuration)
    - [Build the Frontend](#build-the-frontend)
    - [Setup Databases](#setup-databases)
    - [Launch the Backend](#launch-the-backend)
    - [Setup Access to the Backend](#setup-access-to-the-backend)
    - [Serve files for  the Frontend](#serve-files-for--the-frontend)

## Requirements

- Requires node.js 14.4+
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

### Unpack the archive

```shell
tar xvfz ./exais-explorer.preview.1.2.0.tar.gz
cd ./exais-explorer.preview.1.2.0
```

### Review the config files

#### Backend Configuration 

Edit the file  `config_server.json`

This sample JSON expects that the explorer binary runs on the same network node with the blockhain  node.

It exposes REST API and Websocket API through the web server that is served on the port 10001.

Ports for Postgres and Redis must be same as specified in `docker-compose.yaml`

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
      "api": "localhost:42002",
      "grpc": "localhost:42003"
    },
    "db": {
      "driver_name": "postgres",
      "user": "postgres",
      "password": "qwer1234",
      "db_name": "klaatoo-explorer",
      "db_port": 55432,
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
      "api": "localhost:42002",
      "grpc": "localhost:42003"
    },
    "db": {
      "driver_name": "postgres",
      "user": "postgres",
      "password": "qwer1234",
      "db_name": "klaatoo-explorer",
      "db_port": 55432,
      "idle_conn": 2,
      "max_conn": 10
    },
    "redis": {
      "host": "localhost",
      "port": 56379,
      "db": 0
    }
  }
}
```

#### Frontend Configuration

Edit the file  `config_server.json`

```json
{
  "network": {
    "testnet": {
      "name": "Test net",
      "api": "https://explorer.exais.net/",
      "ws": "wss://explorer.exais.net/ws",
      "bech32": {
        "account": "exa",
        "validator": "vl",
        "operator": "op"
      },
      "networkToken": {
        "base": "nexa",
        "symbol": "EXA",
        "precision": 9,
        "description": "ExaisNet Platform Coin"
      }
    }
  }
}
```

### Build the Frontend

Before the next step please install:
- Node.js version 14.4+.
- `yarn` package  manager.

Launch the following code in `exais-explorer.preview.1.2.0` directory:
```bash
./nikto-explorer client  build
```

### Setup Databases

The recommended way  to setup Postgres and Redis is to use docker images.
The file `docker-compose.yaml` provides details on their  initial configuration.

1. Launch  the databases:
   ```shell
   docker-compose up
   ```

2. Setup the schema.
   ```shell
   ./nikto-explorer server setup
   ```

### Launch the Backend

```shell
klaatoo-explorer server run
```

### Setup Access to the Backend

Frontend uses URLs from `config_client.json` to access information. We recommend using `nginx`
reverse proxy to wrap backend's  ports into URLs.

Here is a sample NGINX configuration for the explorer. It assumes that the backend serves on port `10001` and  the frontend server
runs on the port `10000` (see below).

```nginx
server {
  listen 443 ssl;

  server_name explorer.exais.net;
  add_header Access-Control-Allow-Origin *;

  location / {
    proxy_pass http://localhost:10000;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
  }
  location /api/ {
    # Handle preflight requests
    if ($request_method = 'OPTIONS') {
      add_header 'Access-Control-Allow-Origin' '*';
      add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS';
      
      # Custom headers and headers various browsers *should* be OK with but aren't
      add_header 'Access-Control-Allow-Headers' 'DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range';
      
      # Tell client that this pre-flight info is valid for 20 days
      add_header 'Access-Control-Max-Age' 1728000;
      add_header 'Content-Type' 'text/plain; charset=utf-8';
      add_header 'Content-Length' 0;
      return 204;
    }
    if ($request_method = 'POST') {
      add_header 'Access-Control-Allow-Origin' '*' always;
      add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS' always;
      add_header 'Access-Control-Allow-Headers' 'DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range' always;
      add_header 'Access-Control-Expose-Headers' 'Content-Length,Content-Range' always;
    }
    if ($request_method = 'GET') {
      add_header 'Access-Control-Allow-Origin' '*' always;
      add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS' always;
      add_header 'Access-Control-Allow-Headers' 'DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range' always;
      add_header 'Access-Control-Expose-Headers' 'Content-Length,Content-Range' always;
    }
    proxy_pass http://192.168.219.203:10001;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
  }
  location /ws {
    proxy_pass http://192.168.219.203:10001;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "Upgrade";
  }
  ssl_certificate /etc/letsencrypt/live/explorer.exais.net/fullchain.pem; # managed by Certbot
  ssl_certificate_key /etc/letsencrypt/live/explorer.exais.net/privkey.pem; # managed by Certbot
}
```
### Serve files for  the Frontend

```shell
./nikto-explorer client serve --port 10000
```
