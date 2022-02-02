FROM golang:1.13.6 as builder

WORKDIR /elrond
COPY . .
# Proxy
WORKDIR /elrond/cmd/proxy
RUN go build

# ===== SECOND STAGE ======
FROM ubuntu:18.04
COPY --from=builder /elrond/cmd/proxy /elrond/cmd/proxy

WORKDIR /elrond/cmd/proxy/
EXPOSE 8079
ENTRYPOINT ["./proxy"]
