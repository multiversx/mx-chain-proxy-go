# elrond-proxy

The **Elrond Proxy** acts as an entry point into the Elrond Network. 

![Elrond Proxy - Architectural Overview](assets/overview.png "Elrond Proxy - Architectural Overview")

For more details, go to [docs.elrond.com](https://docs.elrond.com/sdk-and-tools/proxy/).

## Rest API endpoints

# V1.0

### address

- `/v1.0/address/:address`         (GET) --> returns the account's data in JSON format for the given :address.
- `/v1.0/address/:address/balance` (GET) --> returns the balance of a given :address.
- `/v1.0/address/:address/nonce`   (GET) --> returns the nonce of an :address.
- `/v1.0/address/:address/shard`   (GET) --> returns the shard of an :address based on current proxy's configuration.
- `/v1.0/address/:address/keys `   (GET) --> returns the key-value pairs of an :address.
- `/v1.0/address/:address/storage/:key`   (GET) --> returns the value for a given key for an account.
- `/v1.0/address/:address/transactions` (GET) --> returns the transactions stored in indexer for a given :address.
- `/v1.0/address/:address/esdt` (GET) --> returns the account's ESDT tokens list for the given :address.
- `/v1.0/address/:address/esdt/:tokenIdentifier` (GET) --> returns the token data for a given :address and ESDT token, such as balance and properties
- `/v1.0/address/:address/esdtnft/:tokenIdentifier/nonce/:nonce` (GET) --> returns the NFT token data for a given address, token identifier and nonce

### transaction

- `/v1.0/transaction/send`         (POST) --> receives a single transaction in JSON format and forwards it to an observer in the same shard as the sender's shard ID. Returns the transaction's hash if successful or the interceptor error otherwise.
- `/v1.0/transaction/simulate`         (POST) --> same as /transaction/send but does not execute it. will output simulation results
- `/v1.0/transaction/simulate?checkSignature=false`         (POST) --> same as /transaction/send but does not execute it, also the signature of the transaction will not be verified. will output simulation results
- `/v1.0/transaction/send-multiple` (POST) --> receives a bulk of transactions in JSON format and will forward them to observers in the rights shards. Will return the number of transactions which were accepted by the interceptor and forwarded on the p2p topic.
- `/v1.0/transaction/send-user-funds` (POST) --> receives a request containing `address`, `numOfTxs` and `value` and will select a random account from the PEM file in the same shard as the address received. Will return the transaction's hash if successful or the interceptor error otherwise.
- `/v1.0/transaction/cost`         (POST) --> receives a single transaction in JSON format and returns it's cost
- `/v1.0/transaction/:txHash` (GET) --> returns the transaction which corresponds to the hash
- `/v1.0/transaction/:txHash?withResults=true` (GET) --> returns the transaction and results which correspond to the hash
- `/v1.0/transaction/:txHash?sender=senderAddress` (GET) --> returns the transaction which corresponds to the hash (faster because will ask for transaction from the observer which is in the shard in which the address is part).
- `/v1.0/transaction/:txHash?sender=senderAddress&withResults=true` (GET) --> returns the transaction and results which correspond to the hash (faster because will ask for transaction from observer which is in the shard in which the address is part)
- `/v1.0/transaction/:txHash/status` (GET) --> returns the status of the transaction which corresponds to the hash
- `/v1.0/transaction/:txHash/status?sender=senderAddress` (GET) --> returns the status of the transaction which corresponds to the hash (faster because will ask for transaction status from the observer which is in the shard in which the address is part).

### vm-values

- `/v1.0/vm-values/hex`            (POST) --> receives a VM Request (`scAddress` string, `funcName` string and `args` []string) and returns the result of the VM Query in hex encoded string format
- `/v1.0/vm-values/string`         (POST) --> receives a VM Request (`scAddress` string, `funcName` string and `args` []string) and returns the result of the VM Query in string format
- `/v1.0/vm-values/int`            (POST) --> receives a VM Request (`scAddress` string, `funcName` string and `args` []string) and returns the result of the VM Query in integer format
- `/v1.0/vm-values/query`          (POST) --> receives a VM Request (`scAddress` string, `funcName` string and `args` []string) and returns the result of the VM Query

### network

- `/v1.0/network/status/:shard`      (GET) --> returns the status metrics from an observer in the given shard
- `/v1.0/network/config`             (GET) --> returns the configuration of the network from any observer
- `/v1.0/network/economics`          (GET) --> returns the economics data metric from the last epoch
- `/v1.0/network/esdts`              (GET) --> returns the names of all the issued ESDTs
- `/v1.0/network/direct-staked-info` (GET) --> returns the list of direct staked values
- `/v1.0/network/delegated-info`     (GET) --> returns the list of delegated values
- `/v1.0/network/enable-epochs`      (GET) --> returns the activation epochs metric
### node

- `/v1.0/node/heartbeatstatus`     (GET) --> returns the heartbeat data from an observer from any shard. Has a cache to avoid many requests

### validator

- `/v1.0/validator/statistics`     (GET) --> returns the validator statistics data from an observer from any shard. Has a cache to avoid many requests

### block

- `/v1.0/block/:shardID/by-nonce/:nonce`    (GET) --> returns a block by nonce
- `/v1.0/block/:shardID/by-nonce/:nonce?withTxs=true`    (GET) --> returns a block by nonce, with transactions included
- `/v1.0/block/:shardID/by-hash/:hash`    (GET) --> returns a block by hash
- `/v1.0/block/:shardID/by-hash/:hash?withTxs=true`    (GET) --> returns a block by hash, with transactions included

### block-atlas

- `/v1.0/block-atlas/:shard/:nonce`   (GET) --> returns a block by nonce, as required by Block Atlas


### hyperblock

- `/v1.0/hyperblock/by-nonce/:nonce`  (GET) --> returns a hyperblock by nonce, with transactions included
- `/v1.0/hyperblock/by-hash/:hash`    (GET) --> returns a hyperblock by hash, with transactions included

### dns
- `/v1.0/dns/all` (GET) --> returns all DNS addresses in a sorted manner
- `/v1.0/dns/username/:username` (GET) --> returns the DNS address corresponding to the provided username

# V_next

This serves as a placeholder for further versions in order to provide a real use-case example of how performing
CRUD operations over endpoints works.
What is different from `v1_0`:
- `/v_next/address/:address/shard` is updated and returns a hardcoded `37` as shard ID
- `/v_next/address/:address/new-endpoint` is added a returns a hardcoded `"test"` as data field
- `/v_next/address/:address/nonce` is removed

The rest of endpoints remain the same.