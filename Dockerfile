FROM golang:1.9.2-alpine3.6
MAINTAINER Toomore Chiang <toomore0929@gmail.com>

WORKDIR /go/src/github.com/toomore/mailbox/

ADD ./campaign ./campaign
ADD ./cmd/campaign.go ./cmd/campaign.go
ADD ./cmd/gendoc.go ./cmd/gendoc.go
ADD ./cmd/root.go ./cmd/root.go
ADD ./cmd/send.go ./cmd/send.go
ADD ./cmd/server.go ./cmd/server.go
ADD ./cmd/user.go ./cmd/user.go
ADD ./mails ./mails
ADD ./main.go ./main.go
ADD ./reader ./reader
ADD ./utils ./utils

VOLUME ["/go/bin"]

RUN \
    apk update && apk add gcc git musl-dev && \
    rm -rf /var/cache/apk/* /var/lib/apk/* /etc/apk/cache/* && \
    go get -v ./...
