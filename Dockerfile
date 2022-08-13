FROM ubuntu:22.04

ARG go_version="1.19"

RUN apt update && \
    apt install -y wget iproute2 net-tools iputils-ping && \
    wget https://dl.google.com/go/go${go_version}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${go_version}.linux-amd64.tar.gz && \
    install /usr/local/go/bin/go /usr/local/bin

WORKDIR /usr/src/app
COPY . .
RUN go mod download
