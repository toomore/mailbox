FROM alpine:latest
MAINTAINER Toomore Chiang <toomore0929@gmail.com>

WORKDIR /bin

ADD ./mailbox_bin ./

RUN apk update && apk add ca-certificates && \
    rm -rf /var/cache/apk/* /var/lib/apk/* /etc/apk/cache/*
