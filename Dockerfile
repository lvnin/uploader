FROM golang:alpine as builder

WORKDIR /go/src/uploader
COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && go build -o server .

FROM alpine:latest

LABEL MAINTAINER="ninlyu.dev@outlook.com"

WORKDIR /go/src/uploader

COPY --from=0 /go/src/uploader ./
COPY --from=0 /go/src/uploader/config.docker.yaml ./

ENV GIN_MODE=release
ENTRYPOINT ./server -c config.docker.yaml
