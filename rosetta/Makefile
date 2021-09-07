GO_PACKAGES=./services/... ./provider/...
TEST_SCRIPT=go test ${GO_PACKAGES}
CURRENT_DIRECTORY=$(shell pwd)
LOCAL_DOCKER_TAG="rosetta-proxy-local:v1.0.0"
OFFLINE_ROSETTA_SERVER_PORT=1234

build-docker-image:
	cd ${CURRENT_DIRECTORY}/..; \
		docker image build . -t elrondnetwork/${LOCAL_DOCKER_TAG} -f ./rosetta/Dockerfile; \
		cd ${CURRENT_DIRECTORY}
	echo ${LOCAL_DOCKER_TAG} > ~/.proxyDockerTag

run-offline:
	./rosetta.sh offline ${OFFLINE_ROSETTA_SERVER_PORT}

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
