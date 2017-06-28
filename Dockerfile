FROM golang:alpine
MAINTAINER Toomore Chiang <toomore0929@gmail.com>

WORKDIR /go/src/github.com/toomore/mailbox/

ADD ./campaign ./campaign
ADD ./cmd/mailbox_sender/main.go ./cmd/mailbox_sender/main.go
ADD ./mailbox/cmd ./mailbox/cmd
ADD ./mailbox/main.go ./mailbox/main.go
ADD ./mails ./mails
ADD ./reader ./reader
ADD ./utils ./utils

VOLUME ["/go/bin"]

RUN \
    apk update && apk add gcc git musl-dev && \
    rm -rf /var/cache/apk/* /var/lib/apk/* /etc/apk/cache/* && \
    go get -v ./...
