## Overview
`rosetta-elrond` provides a reference implementation of the Rosetta API for
Elrond blockchain in Golang. If you haven't heard of the Rosetta API, you can find more
information [here](https://rosetta-api.org).


## Features
* Rosetta API implementation (both Data API and Construction API)
* UTXO cache for all accounts (accessible using `/account/balance`)
* Stateless, offline, curve-based transaction construction from any Bech32 Address

## Usage
As specified in the [Rosetta API Principles](https://www.rosetta-api.org/docs/automated_deployment.html),
all Rosetta implementations must be deployable via Docker and support running via either an
[`online` or `offline` mode](https://www.rosetta-api.org/docs/node_deployment.html#multiple-modes).

**YOU MUST INSTALL DOCKER FOR THE FOLLOWING INSTRUCTIONS TO WORK. YOU CAN DOWNLOAD
DOCKER [HERE](https://www.docker.com/get-started).**

### Install
Running the following commands will create a Docker image called `rosetta-elrond:latest`.


## Future Work
* [Rosetta API `/mempool`](https://www.rosetta-api.org/docs/MempoolApi.html)
* [Rosetta API `/mempool/transaction`](https://www.rosetta-api.org/docs/MempoolApi.html#mempooltransaction) implementation
