GO_PACKAGES=./services/... ./client/...
TEST_SCRIPT=go test ${GO_PACKAGES}

run-mainnet:
	./rosetta.sh mainnet

run-testnet:
	./rosetta.sh testnet

stop:
	./rosetta.sh stop

test:
	${TEST_SCRIPT}

coverage-local:
	${TEST_SCRIPT} -cover
