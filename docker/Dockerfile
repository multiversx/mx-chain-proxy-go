FROM golang:1.22.4 as builder

WORKDIR /mx-chain-proxy-go
COPY . .

# Proxy
WORKDIR /mx-chain-proxy-go/cmd/proxy
RUN go build -ldflags="-X main.appVersion=$(git describe --tags --long --dirty) -X main.commitID=$(git rev-parse HEAD)"

# ===== SECOND STAGE ======
FROM ubuntu:22.04
RUN apt-get update -y && apt-get upgrade -y
COPY --from=builder /mx-chain-proxy-go/cmd/proxy /mx-chain-proxy-go/cmd/proxy

WORKDIR /mx-chain-proxy-go/cmd/proxy/
EXPOSE 8080
ENTRYPOINT ["./proxy"]
