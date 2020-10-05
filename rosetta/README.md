## Overview

This is the reference implementation of the [Rosetta API](https://rosetta-api.org) for Elrond.

The Rosetta API has been implemented as an extension of the [Elrond Proxy](https://github.com/ElrondNetwork/elrond-proxy-go).

The implementation is supported by an [Observing Squad](https://docs.elrond.com/observing-squad), where the Proxy starts as gateway that resolves the impedance mismatch between the Elrond API (exposed the the Observer Nodes) and the Rosetta API.


Note: An **Observing Squad** is defined as a set of `N` Observer Nodes (one for each Shard, including the Metachain) plus the Elrond Proxy instance (which connects to these Observers and delegates requests towards them). Currently the Elrond Mainnet has 3 Shards, plus the Metachain. Therefore, the Observing Squad is composed of 4 Observers and one Proxy instance.


One can set up a Rosetta-enabled Elrond Observing Squad via the provided scripts and Makefile - which use `docker` and `docker-compose` under the hood.

## Features

* Rosetta API implementation (both Data API and Construction API)
* Stateless, offline, curve-based transaction construction from any Bech32 Address

## Prerequisites

You need to install Docker and Docker Compose.

For example, on Ubuntu:

```
sudo apt-get docker
sudo apt-get docker-compose
```

## Build

In order to build the `rosetta-client` docker image, run the following command:

```
make build-docker-image
```

Under the hood, this command runs `docker build` against the Dockerfile `elrond-proxy`.

## Start

Running the commands below will start a Rosetta-enabled Observing Squad (4 observer nodes, plus the Proxy). The API will be available at the following address: `http://10.0.0.2:8079`.

### Testnet

```
make run-testnet
```

### Mainnet

```
make run-mainnet
```

## Stop

In order to stop the Observing Squad, run the command:

```
make stop
```

## System Requirements

The system requirements for an Observing Squad are listed [here](https://docs.elrond.com/observing-squad#system-requirements).

## Testing with `rosetta-cli`

In order to validate the Elrond implementation of the Rosetta API, [install `rosetta-cli`](https://github.com/coinbase/rosetta-cli#install) and run one of the following commands:

* `rosetta-cli check:data --configuration-file rosetta-cli-conf/elrond_testnet.json`
* `rosetta-cli check:construction --configuration-file rosetta-cli-conf/elrond_testnet.json`
* `rosetta-cli check:data --configuration-file rosetta-cli-conf/elrond_mainnet.json`
* `rosetta-cli check:construction --configuration-file rosetta-cli-conf/elrond_mainnet.json`

## Future Work

* [Rosetta API `/mempool`](https://www.rosetta-api.org/docs/MempoolApi.html)
* [Rosetta API `/mempool/transaction`](https://www.rosetta-api.org/docs/MempoolApi.html#mempooltransaction)
