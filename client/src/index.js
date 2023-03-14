import App from './App.svelte';
import * as config from "./config.json";
import i18n from './i18n';
import {stats, websocket} from 'js';
import {getNetwork, setNetwork} from "js/network";
import * as style from './styles/style.css'

console.log("index.js")

function init() {
    
    // load config
    window.config = config;

    const network = getNetwork();
    setNetwork(network);

    // set i18n data
    i18n();

    // create connection with websocket server
    const conn = new websocket();
    window.klaatoo = conn;

    // connect to server
    conn.connect(() => {
        conn.subscribe({
            method: 'stats.now',
            params: [],
            id: conn.generateRequestId(),
            success: (data) => {
                stats.set(data);
            },
            error: (error) => {
                console.error(error);
            },
        });
    });
}

init();

export default new App({
    target: document.body,
});
