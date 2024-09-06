FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o  storage-accessor .

FROM alpine
COPY --from=builder /app/storage-accessor /bin/storage-accessor

ENTRYPOINT ["/bin/storage-accessor"]
