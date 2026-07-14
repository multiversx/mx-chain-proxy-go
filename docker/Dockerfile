FROM golang:1.26.2 AS builder
LABEL maintainer="multiversx"

WORKDIR /mx-chain-proxy-go
COPY . .

# Proxy
WORKDIR /mx-chain-proxy-go/cmd/proxy
RUN go build -ldflags="-X main.appVersion=$(git describe --tags --long --dirty) -X main.commitID=$(git rev-parse HEAD)"

# ===== SECOND STAGE ======
FROM ubuntu:22.04
RUN apt-get update -y && apt-get upgrade -y

RUN useradd --create-home --uid 10001 --shell /usr/sbin/nologin appuser

COPY --from=builder /mx-chain-proxy-go/cmd/proxy /mx-chain-proxy-go/cmd/proxy

RUN chown -R appuser:appuser /mx-chain-proxy-go/cmd/proxy/

WORKDIR /mx-chain-proxy-go/cmd/proxy/
USER appuser
EXPOSE 8080
ENTRYPOINT ["./proxy"]
