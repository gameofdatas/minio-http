FROM golang:alpine AS builder
WORKDIR /go/src/github.com/k8-proxy/storage-accessor
COPY . .
RUN env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o  storage-accessor .

FROM alpine
COPY --from=builder /go/src/github.com/k8-proxy/storage-accessor/storage-accessor /bin/storage-accessor

ENTRYPOINT ["/bin/storage-accessor"]