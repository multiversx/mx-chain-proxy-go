FROM golang:1.13.6 as builder

WORKDIR /elrond
COPY . .
# Proxy
WORKDIR /elrond/cmd/proxy
RUN go build


# ===== SECOND STAGE ======
FROM ubuntu:18.04
COPY --from=builder /elrond/cmd/proxy /elrond/cmd/proxy
# COPY config.toml file from rosetta folder to proxy config folder
COPY --from=builder /elrond/rosetta/config.toml /elrond/cmd/proxy/config/
COPY --from=builder /elrond/rosetta/offline_config_devnet.toml /elrond/cmd/proxy/config/
COPY --from=builder /elrond/rosetta/offline_config_mainnet.toml /elrond/cmd/proxy/config/

WORKDIR /elrond/cmd/proxy/
EXPOSE 8079
ENTRYPOINT ["./proxy", "--rosetta"]
