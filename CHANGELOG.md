# Changelog

## v0.3.0
> Deleted tailwind, daisy ui dependencies
> 
> Reimplemented and redesigned client by @czmotion
---

### Client

* feat : reimplemented & redesigned client by @czmotion

### Server 

* fast-sync : update txs across multiple goroutines.
* feat : tokens api with paginate
* feat : validators api 
* fix : panic when token amount is over int64


## v0.2.0

### Client
> Replaced tailwind plugin to daisy ui
---

* feat : added command to fetch css data.
* feat : added tokens list view.
* feat : added network status label.
* stats : added txs count statistic in dashboard
* fix : Integrate colors to semantic color

### Server

* fix : concurrency panic when access cache data on memory.

## v0.1.3

### Client
* fix : Duplicate request id [#16]("https://github.com/hessegg/klaatoo-explorer/issues/16")
* fix : Replaced average time with average best time. [#18]("https://github.com/hessegg/klaatoo-explorer/pull/18)

### Server
* sync : Insert blocks and transactions in a single db transaction.
* fast-sync : sync blocks across multiple goroutines.
* event : fix api to paginate data.
* event : fix caching data on event handler.
* service : manage service with ServiceManager
* handler : manage handler with HandlerManager
* swagger : generated swagger document

## v0.1.2

### Client
* websocket : fix registering duplicate events.
* websocket : try to reconnect when connection closed.
* cmd : support setup database command
* websocket : handle websocket response as ajax style.
* i18n : added more translation data.

### Server
* common: handle goroutines with context
* api: integrated into websocket server.
* websocket : `WSResponse` includes response code.


## v0.1.1

### Client
* websocket : generate random request id and added to ws request.

### Server
* websocket : added request id to `wsEvent` and uses in event key.


## v0.1.0

### Client 
* i18n : locale translation
* websocket : replaced RESTApi connection to websocket.

### Server
* websocket : modified random string session id to uuid
* websocket : Added `cache` option to bypass common db query.
* websocket : replaced RESTApi connection to websocket.
* sync : recover goroutine when sync process panic.

## v0.0.1

### Client
* Explore recent blocks and transactions
* Explore block details
* Explore transaction details
* Explore address details

### Server
* sync: block, transaction, validator synchronization


* api: RESTapi for query block, transaction and validator data
