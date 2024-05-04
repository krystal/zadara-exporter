FROM alpine:3.19
# RUN apk update && apk add --no-cache ca-certificates curl
ENTRYPOINT [ "/zadara-exporter" ]
COPY zadara-exporter /
