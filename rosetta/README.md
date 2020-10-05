## Overview

`rosetta-elrond` provides a reference implementation of the Rosetta API for
Elrond blockchain in Golang. If you haven't heard of the Rosetta API, you can find more
information [here](https://rosetta-api.org).

Elrond rosetta client is an extension to our elrond proxy implementation that is a gateway to elrond-blockchain

In order to work correctly we must have 4 observers nodes, a node for every shard  (our mainnet has 3 shards 
plus a metachain shard) that will provide information about every shard.
To can simplify this configuration we have a `docker compose configuration` that will start all 4 observers nodes and rosetta client. 

## Features

* Rosetta API implementation (both Data API and Construction API)
* Stateless, offline, curve-based transaction construction from any Bech32 Address


*YOU MUST INSTALL DOCKER FOR THE FOLLOWING INSTRUCTIONS TO WORK. YOU CAN DOWNLOAD
DOCKER [HERE](https://www.docker.com/get-started).*

*YOU ALSO HAVE TO INSTALL DOCKER COMPOSE*

```
sudo apt-get docker-compose
```


## START

Running the following commands will start an elrond rosetta-client with observing-squad (4 observer nodes)

A rosetta-client will start at address: `http://10.0.0.2:8079`

### Testnet

```

make run-testnet
```
### Mainnet

```
make run-mainnet
```

### To stop `elrond-roseta-client` and `observing-squad`

```
make stop
```


##System Requirements

Elrond rosetta client was tested on a machine with 16GBs of RAM and a 6 core CPU

##Testing with `rosetta-cli`

To validate `rosetta-elrond`, [install `rosetta-cli`](https://github.com/coinbase/rosetta-cli#install)
and run one of the following commands:
* `rosetta-cli check:data --configuration-file rosetta-cli-conf/elrond_testnet.json`
* `rosetta-cli check:construction --configuration-file rosetta-cli-conf/elrond_testnet.json`
* `rosetta-cli check:data --configuration-file rosetta-cli-conf/elrond_mainnet.json`
* `rosetta-cli check:construction --configuration-file rosetta-cli-conf/elrond_mainnet.json`

## Future Work

* [Rosetta API `/mempool`](https://www.rosetta-api.org/docs/MempoolApi.html)
* [Rosetta API `/mempool/transaction`](https://www.rosetta-api.org/docs/MempoolApi.html#mempooltransaction) implementation
