FROM golang:1.17 as builder

WORKDIR /multiversx
COPY . .

# Proxy
WORKDIR /multiversx/cmd/proxy
RUN go build -ldflags="-X main.appVersion=$(git describe --tags --long --dirty) -X main.commitID=$(git rev-parse HEAD)"

# ===== SECOND STAGE ======
FROM ubuntu:18.04
COPY --from=builder /multiversx/cmd/proxy /multiversx/cmd/proxy

WORKDIR /multiversx/cmd/proxy/
EXPOSE 8079
ENTRYPOINT ["./proxy"]
