# elrond-proxy

**Rest API endpoints:** 

-*address* group
- `/address/:address`         (GET) --> returns the account's data in JSON format for the given :address.
- `/address/:address/balance` (GET) --> returns the balance of a given :address.
- `/address/:address/nonce`   (GET) --> returns the nonce of an :address.
- `/address/:address/storage/:key`   (GET) --> returns the value for a given key for an account.
- `/address/:address/transactions` (GET) --> returns the transactions stored in indexer for a given :address.

-*transaction* group
- `/transaction/send`         (POST) --> receives a single transaction in JSON format and forwards it to an observer in the same shard as the sender's shard ID. Returns the transaction's hash if successful or the interceptor error otherwise.
- `/transaction/send-multiple` (POST) --> receives a bulk of transactions in JSON format and will forward them to observers in the rights shards. Will return the number of transactions which were accepted by the interceptor and forwarded on the p2p topic.
- `/transaction/send-user-funds` (POST) --> receives a request containing `address`, `numOfTxs` and `value` and will select a random account from the PEM file in the same shard as the address received. Will return the transaction's hash if successful or the interceptor error otherwise.
- `/transaction/cost`         (POST) --> receives a single transaction in JSON format and returns it's cost
- `/transaction/:txHash` (GET) --> returns the transaction which corresponds to the hash
- `/transaction/:txHash?sender=senderAddress` (GET) --> returns the transaction which corresponds to the hash (faster because will ask for transaction from observer which is in the shard in which the address is part)
- `/transaction/:txHash/status` (GET) --> returns the status of the transaction which corresponds to the hash

-*vm-values* group
- `/vm-values/hex`            (POST) --> receives a VM Request (`scAddress` string, `funcName` string and `args` []string) and returns the result of the VM Query in hex encoded string format
- `/vm-values/string`         (POST) --> receives a VM Request (`scAddress` string, `funcName` string and `args` []string) and returns the result of the VM Query in string format
- `/vm-values/int`            (POST) --> receives a VM Request (`scAddress` string, `funcName` string and `args` []string) and returns the result of the VM Query in integer format
- `/vm-values/query`          (POST) --> receives a VM Request (`scAddress` string, `funcName` string and `args` []string) and returns the result of the VM Query

-*network* group
- `/network/status/:shard`    (GET) --> returns the status metrics from an observer in the given shard
- `/network/config`           (GET) --> returns the configuration of the network from any observer

*transaction* group
- `/node/heartbeatstatus`     (GET) --> returns the heartbeat data from an observer from any shard. Has a cache to avoid many requests

*validator* group
- `/validator/statistics`     (GET) --> returns the validator statistics data from an observer from any shard. Has a cache to avoid many requests

*block* group
- `/block/:shardID/:nonce`    (GET) --> will search in the indexer for the block in the given shard and at the given nonce and returns it
