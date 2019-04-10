FROM alpine:3.9
MAINTAINER Toomore Chiang <toomore0929@gmail.com>

WORKDIR /bin

ADD ./mailbox_bin ./

RUN apk update && apk add ca-certificates bash-completion && \
    cd /usr/share/bash-completion/completions && mailbox doc -b && \
    rm -rf /var/cache/apk/* /var/lib/apk/* /etc/apk/cache/*

WORKDIR /workdir

CMD ["/bin/bash", "-l"]
