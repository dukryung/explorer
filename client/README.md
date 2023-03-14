![Language](https://img.shields.io/badge/Language-Svelte%203.4-red.svg)
![UI](https://img.shields.io/badge/Node.js-%20up_to_17-blue)
![Xcode](https://img.shields.io/badge/Webpack-5.37.1-green)
![Version](https://img.shields.io/badge/Version-1.0.0-purple)


# Nextor Explorer

## Table of Contents

- [Requirements](#requirements)
- [Quick Start](#quick-start)
- [Layout](#layout)
- [Deploy](#deploy)

## Requirements

- svelte 3.46.4+
- node.js up to Version 17 (17 not included)
- webpack 5.37.1

## Quick Start


1. Open terminal and install the dependencies...

```shell
cd explorer_project_folder
npm install or yarn install,
npm start or yarn start
```
Great! You can develop the project in your localhost: 8080

2. Setup your configuration

`src/config.json`

```json
{
  "network": {
    "mainnet": {
      "name": "Hostnet",
      "api": "https://explorer.exais.net/rest",
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
        "description": "ExaisNet Platorm Coin"
      }
    },
    "testnet": {
      "name": "Hostnet",
      "api": "https://explorer.exais.net/rest",
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
        "description": "ExaisNet Platorm Coin"
      }
    },
    "localnet": {
      "name": "Hostnet",
      "api": "https://explorer.exais.net/rest",
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
        "description": "ExaisNet Platorm Coin"
      }
    }
  }
}
```


- `network` add a network to the application by a key like `mainnet` and `testnet`.
  - `name` network name on the client side.
  - `ws` is websocket URL  
  - `bech32` cosmos-sdk bech32 config.
  - `networkToken` default network token information.

## Layout

```tree
├── public
│   ├── images
│   └── index.html
├── scripts
│   ├── fetch.i18n.js 
│   └── google-auth.json  
├── src
│   ├── assets
│   ├── components
│   │   ├── Cards
│   │   ├── DropDown
│   │   ├── Footers 
│   │   ├── Headers 
│   │   ├── Labels 
│   │   ├── Messages
│   │   └── TableItems
│   ├── i18n
│   │   ├── en-Us.json
│   │   ├── index.js
│   │   ├── ko-KR.json
│   │   └── ru-Ru.json
│   ├── js
│   │   ├── index.js
│   │   ├── network.js
│   │   ├── store.js
│   │   ├── token.js
│   │   ├── util.js
│   │   └── websocket.js
│   ├── layouts
│   ├── styles
│   ├── App.svelte
│   ├── config.json
│   └── index.js
├── webpack
│   ├── utils
│   │   ├── environment.js
│   │   ├── index.js
│   │   ├── paths.js
│   │   ├── plugins.js 
│   │   └── sever.js
│   ├── build.js
│   ├── common.js 
│   └── start.js
├── .eslintignore
├── .eslintrc.js
├── package-lock.json
├── package.json
├── postcss.config.js
├── README.md
├── tailwind.config.js
├── yarn.lock
```
A brief explanation of a layout:

* `public` has images that we have used in the project and a root html file.
* `scripts` folder owns files which are used for i18n dictionary JSON file from google-spreadsheet.
* `src` directory contains the meat of the project.
  - `assets` has all icons and assets which is used in the project.
  - `components` folder owns all necessary and reusable components.
  - `i18n` is used for designing our project for use in different locales around the world.
  - `js` has all the core functionality of our Project.
    - `index.js` helps export all files as one file.
    - `network.js` checks network connection in the wallet like testnet or mainnet is connected or not.
    - `store.js` is a writable store that has stats, accessState, networkConnection, selectedAccount, and so on. you can use this value any place of the project directly.
    - `token.js` is a Token class in which you can initialize a token based on config file.
    - `util.js` has all utility functions which is used over the project
    - `websocket.js` has a websocket class. there are two types of requests in it. single request or subscription like receiving real-time data.
  - `layouts` has all main UI components including Dashboard and all section components 
  - `styles` has our css files
  - `App.svelte` is a main root svelte file and all routing starts from here
  - `config.json` is a config file that is explained above briefly.
  - `index.js` is a main root js file which is a project that starts to run from here. All set-up stuff is prepared here.
* `webpack` module bundler which makes a large number of files into a single file
* `.eslintignore` to find a .eslintignore file before determining which files to lint.
* `.eslintrc.js` a configuration file for a tool, ESLINT
* `package-lock.json` is automatically generated for any operations where npm modifies either the node_modules tree or package.json
* `package.json` enables npm to start your project, run scripts, install dependencies, publish to the NPM registry, and many other useful tasks.
* `README.md` is a detailed description of the project.
* `tailwind.config.js` is a tailwind config file in which you can set up the configuration.
* `yarn.lock` describes a project's dependency graph: direct dependencies, child dependencies, and so on


## Deploy
Build or serve

```shell 
npm run build

or

yarn build
```

- Dev : [Link](https://explorer.exais.net/)



