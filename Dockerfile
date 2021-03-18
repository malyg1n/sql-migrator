FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY ./ /app

RUN mkdir -p /.cache && chmod 777 /.cache

RUN go get github.com/githubnemo/CompileDaemon

RUN go mod download