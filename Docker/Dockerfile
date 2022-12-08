FROM golang:1.17 as builder

WORKDIR /elrond
COPY . .

# Proxy
WORKDIR /elrond/cmd/proxy
RUN go build -ldflags="-X main.appVersion=$(git describe --tags --long --dirty) -X main.commitID=$(git rev-parse HEAD)"

# ===== SECOND STAGE ======
FROM ubuntu:18.04
COPY --from=builder /elrond/cmd/proxy /elrond/cmd/proxy

WORKDIR /elrond/cmd/proxy/
EXPOSE 8079
ENTRYPOINT ["./proxy"]
