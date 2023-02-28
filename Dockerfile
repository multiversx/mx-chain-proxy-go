FROM golang:alpine as builder

WORKDIR /elrond
COPY . .
# Proxy
WORKDIR /elrond/cmd/proxy
RUN go build

# ===== SECOND STAGE ======
FROM ubuntu:22.04
# RUN apt-get update && apt-get install -y openssl ca-certificates

COPY --from=builder /elrond/cmd/proxy /elrond/cmd/proxy

WORKDIR /elrond/cmd/proxy/
EXPOSE 8079
ENTRYPOINT ["./proxy", "--start-swagger-ui"]
