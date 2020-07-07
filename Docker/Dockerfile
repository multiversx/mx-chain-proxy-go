FROM golang:1.13.6 as builder

WORKDIR /elrond
COPY . .
# Proxy
WORKDIR /elrond/cmd/proxy
RUN go build

# ===== SECOND STAGE ======
FROM ubuntu:18.04
COPY --from=builder /elrond/cmd/proxy /elrond/cmd/proxy
##Copy the proxy config.toml file
#COPY CONFIGS/config.toml /go/elrond-proxy-go/cmd/proxy/config/
#
##Copy over the pems
#COPY NODE-CONFIGS/walletKey.pem /go/elrond-proxy-go/cmd/proxy/config
#COPY NODE-CONFIGS/walletKey.pem /go/elrond-txgen-go/cmd/txgen/config
WORKDIR /elrond/cmd/proxy/
EXPOSE 8079
CMD ["./proxy"]
