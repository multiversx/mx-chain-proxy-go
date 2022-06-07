module github.com/ElrondNetwork/elrond-proxy-go

go 1.15

require (
	github.com/ElrondNetwork/elastic-indexer-go v1.2.30
	github.com/ElrondNetwork/elrond-go v1.3.27
	github.com/ElrondNetwork/elrond-go-core v1.1.15
	github.com/ElrondNetwork/elrond-go-crypto v1.0.1
	github.com/ElrondNetwork/elrond-go-logger v1.0.7
	github.com/ElrondNetwork/elrond-vm-common v1.2.14
	github.com/coinbase/rosetta-sdk-go v0.7.0
	github.com/elastic/go-elasticsearch/v7 v7.12.0
	github.com/gin-contrib/cors v0.0.0-20190301062745-f9e10995c85a
	github.com/gin-contrib/pprof v1.3.0
	github.com/gin-gonic/gin v1.7.6
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	github.com/urfave/cli v1.22.5
	gopkg.in/go-playground/validator.v8 v8.18.2
)

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_2 v1.2.35 => github.com/ElrondNetwork/arwen-wasm-vm v1.2.35

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_3 v1.3.35 => github.com/ElrondNetwork/arwen-wasm-vm v1.3.35

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_4 v1.4.34-rc14 => github.com/ElrondNetwork/arwen-wasm-vm v1.4.34-rc14
