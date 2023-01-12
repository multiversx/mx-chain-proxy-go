FROM golang:1.17.6 as builder

WORKDIR /mx-chain-proxy
COPY . .

# Proxy
WORKDIR /mx-chain-proxy/cmd/proxy
RUN go build

# ===== SECOND STAGE ======
FROM ubuntu:22.04
RUN apt-get update -y && apt-get upgrade -y
COPY --from=builder /mx-chain-proxy/cmd/proxy /mx-chain-proxy/cmd/proxy

WORKDIR /mx-chain-proxy/cmd/proxy/
EXPOSE 8080
ENTRYPOINT ["./proxy"]
