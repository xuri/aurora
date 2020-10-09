FROM golang:alpine AS builder

RUN \
    apk -U upgrade --no-cache && \
    apk add --no-cache build-base git

COPY . /tmp/aurora
RUN \
    cd /tmp/aurora && \
    go build

FROM alpine

COPY --from=builder /tmp/aurora/aurora /usr/bin/
COPY --from=builder /tmp/aurora/aurora.toml /etc/

RUN \
    sed -i "s/127.0.0.1/0.0.0.0/g" /etc/aurora.toml

EXPOSE 3000
ENTRYPOINT ["/usr/bin/aurora", "-c", "/etc/aurora.toml"]
