FROM golang:1.18

WORKDIR /usr/src/app

RUN mkdir -p /dev/net && \
    mknod /dev/net/tun c 10 200 && \
    chmod 666 /dev/net/tun

COPY go.mod ./
RUN go mod download && go mod verify

COPY . .
