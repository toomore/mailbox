#!/usr/bin/env bash
docker run -it --rm  --link mailbox-mariadb:MARIADB              \
           -v $(pwd)/cmd:/cmd                                    \
           ubuntu:16.04 bash

