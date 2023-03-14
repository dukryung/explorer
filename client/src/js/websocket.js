import { networkConnection } from "../js";

export class websocket {
    SINGLE_REQUEST = "SINGLE_REQUEST";
    SUBSCRIBE = "SUBSCRIBE";
    UNSUBSCRIBE = "UNSUBSCRIBE";

    RANDOM_RANGE = 10000;

    constructor(url) {
        if (url === undefined) {
            this.url = network.ws;
        } else {
            this.url = url;
        }
        this.connected = false;
        this.events = [];
    }

    connect(cb) {
        this.socket = new WebSocket(this.url);

        this.socket.onopen = (e) => {
            // console.debug(`connected to ${this.url}`)
            this.setConnected(true);

            if (typeof cb === "function") {
                cb();
            }
        }

        this.socket.onmessage = (e) => {
            const message = JSON.parse(e.data);
            for (const event of this.events) {
                if (message.code !== 0) {
                    if (typeof event.error === 'function') {
                        event.error(message.result)
                    }
                    this.removeEvent(event);
                    return
                }

                if (event.method === message.method && event.id === message.id) {
                    if (typeof event.success === 'function') {
                        event.success(message.result);
                    }

                    switch (event.type) {
                        case this.SINGLE_REQUEST:
                            this.removeEvent(event);
                            break;
                        default:
                    }
                }
            }
        }

        this.socket.onerror = (e) => {
        }

        this.socket.onclose = (e) => {
            // console.debug('disconnected');
            this.setConnected(false);
            this.reconnect();
        }
    }

    setConnected(connected) {
        this.connected = connected;
        networkConnection.set(connected)
    }

    disconnect() {
        if (this.connected) {
            this.socket.close();
        }
    }

    reconnect() {
        if (!this.connected) {
            // console.debug(`try to reconnect websocket server ${this.url}`)
            this.connect(() => {
                let e = this.events;
                for (const event of e) {
                    this.request(event)
                }
            })
        }
    }

    /**
     * Subscribe websocket event
     * @param method method
     * @param params []any
     * @param id request id
     * @param success success callback
     * @param error error callback
     */
    subscribe({ method, params, id, success, error }) {
        this.request({
            type: this.SUBSCRIBE,
            method: method,
            params: params,
            id: id,
            success: success,
            error: error,
        })
    }

    /**
     * Single request websocket event
     * @param method method
     * @param params []any
     * @param id request id
     * @param success success callback
     * @param error error callback
     */
    singleRequest({ method, params, id, success, error }) {
        this.request({
            type: this.SINGLE_REQUEST,
            method: method,
            params: params,
            id: id,
            success: success,
            error: error,
        })
    }

    /**
     * Unsubscribe websocket event
     * @param method method
     * @param params []any
     * @param id request id
     */
    unsubscribe({ method, params, id }) {
        this.request({
            type: this.UNSUBSCRIBE,
            method: method,
            params: params,
            id: id,
        })
    }

    request({ type, method, params, id, success, error }) {
        if (typeof type === 'undefined' || typeof method === 'undefined') {
            if (typeof error === 'function') {
                error("invalid request")
            } else {
                console.error("invalid request")
            }
            return
        }

        // console.debug("REQUEST", type, method, params, id)

        this.addEvent({ type, method, params, id, success, error })
    }

    send(event) {
        this.waitConnection(
            () => {
                this.socket.send(JSON.stringify(event));
            },
            500,
            event,
        )
    }

    waitConnection(cb, iv, event) {
        if (this.socket.readyState === 1) {
            // console.debug("waiting connection: cb");
            cb();
        } else {
            setTimeout(() => {
                if (this.events.find(e => e.id === event.id) !== 'undefined'
                    || event.type === this.UNSUBSCRIBE) {
                    this.waitConnection(cb, iv, event);
                } else {
                    // console.debug("waiting connection: exit :event removed");
                }
            }, iv)
        }
    }

    addEvent(event) {
        switch (event.type) {
            case this.UNSUBSCRIBE:
                this.removeEvent(event);
                break;
            default:
                if (typeof this.events.find(e => e.id === event.id) === 'undefined') {
                    this.events.push(event);
                } else {
                    // console.debug("add event: already registered event", event.method, event.id)
                }
                break;
        }
        this.send(event);
    }

    removeEvent(event) {
        // console.log("remove event ", event.id)
        this.events = this.events.filter(e => e !== event);
    }

    /**
     * Generate request id. range 0~10000
     * @returns {number}
     */
    generateRequestId() {
        return Math.floor(Math.random() * this.RANDOM_RANGE)
    }
}