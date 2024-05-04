FROM golang:1.22-alpine3.19 AS builder
WORKDIR /app
ENV CGO_ENABLED 0
COPY . .
RUN 	go build -o zadara-exporter -ldflags "-w -s \
-X main.version=$(shell git describe --tags --exact 2>/dev/null) \
-X main.commit=$($(shell git rev-parse HEAD 2>/dev/null))"

FROM alpine:3.19
WORKDIR /app
RUN apk add --no-cache \
    ca-certificates \
    curl
COPY --from=builder /app/zadara-exporter /app/zadara-exporter
CMD ["./zadara-exporter", "server"]
