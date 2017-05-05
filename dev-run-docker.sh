#!/usr/bin/env bash
docker run -it --rm  --link mailbox-mariadb:MARIADB              \
           -v $(pwd)/cmd:/cmd                                    \
           -v $(pwd)/campaign:/go/src/github.com/toomore/mailbox/campaign      \
           -v $(pwd)/utils:/go/src/github.com/toomore/mailbox/utils            \
           golang:1.8.1 bash

