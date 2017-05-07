#!/bin/bash

docker run -it --rm -v $(pwd)/mailbox_bin:/mailbox_bin toomore/mailbox:base    \
    sh -c "go get -v ./cmd/...;
           cp /go/bin/* /mailbox_bin;"

docker build -t toomore/mailbox:cmd ./

sudo rm -rf ./mailbox_bin
