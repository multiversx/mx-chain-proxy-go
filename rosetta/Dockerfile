FROM golang:1.13.6 as builder

WORKDIR /elrond
COPY . .
# Proxy
WORKDIR /elrond/cmd/proxy
RUN go build
# $(ls /go/pkg/mod/github.com/!elrond!network/ | grep arwen-wasm-vm | sed 's/.* //' | head -1) this command is required  to be able to exact the full name of folder that starting with 'arwen-wasm-wm'
RUN cp /go/pkg/mod/github.com/!elrond!network/$(ls /go/pkg/mod/github.com/!elrond!network/ | grep arwen-wasm-vm | sed 's/.* //' | head -1)/wasmer/libwasmer_linux_amd64.so /lib/libwasmer_linux_amd64.so

# ===== SECOND STAGE ======
FROM ubuntu:18.04
COPY --from=builder /elrond/cmd/proxy /elrond/cmd/proxy
# COPY config.toml file from rosetta folder to proxy config folder
COPY --from=builder /elrond/rosetta/config.toml /elrond/cmd/proxy/config/
COPY --from=builder "/lib/libwasmer_linux_amd64.so" "/lib/libwasmer_linux_amd64.so"

WORKDIR /elrond/cmd/proxy/
EXPOSE 8079
ENTRYPOINT ["./proxy", "--rosetta"]
