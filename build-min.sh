#!/bin/bash

docker run -it --rm -v $(pwd)/mailbox_bin:/mailbox_bin toomore/mailbox:base    \
    sh -c "GO111MODULE=on go get -v ./...;
           cp /go/bin/* /mailbox_bin;"

docker build -t toomore/mailbox:cmd -f ./Dockerfile-min ./

sudo rm -rf ./mailbox_bin
