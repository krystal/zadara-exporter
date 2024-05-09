FROM golang:1.22.2-alpine as builder

WORKDIR /work

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -mod=readonly -a -o zadara-exporter .

FROM alpine:3 as certs

RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

WORKDIR /
COPY --from=builder /zadara-exporter /
COPY ./docker-entrypoint.sh /
ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["/zadara-exporter", "server"]