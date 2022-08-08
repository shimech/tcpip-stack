FROM golang:1.19

WORKDIR /usr/src/app
COPY . .
RUN go mod download
