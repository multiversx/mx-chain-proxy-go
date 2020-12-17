### General
Api testing tool that is written in typescript and uses `erdjs`.
It runs in browser and allows testing the HTTP Rest API endpoints of proxy by just clicking
some buttons. It also has a set of checks for each endpoint, but the request and the response
are also displayed, so additional manual checks can be done. 
For now, it tests almost all endpoints in proxy (block group is missing for example).
The code should allow easy integration for new endpoints, versions and checks.

### Prerequisites

`browserify` and `http-server` are required. They can be installed as follows:

    npm install --global browserify
    npm install --global http-server

### Install & running

    npm install (first time only)
    npm run compile
    npm run serve
    open localhost:7777 in browser.
