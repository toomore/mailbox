#!/usr/bin/env bash
docker run -d --name mailbox-prod-web                                       \
           --link mailbox-mariadb-prod:MARIADB                              \
           --log-opt max-size=64m                                           \
           --log-opt max-file=1                                             \
           --restart=always                                                 \
           -p 127.0.0.1:8801:8801                                           \
           toomore/mailbox:cmd mailbox_server

docker run -d --name mailbox-prod-web-2                                     \
           --link mailbox-mariadb-prod:MARIADB                              \
           --log-opt max-size=64m                                           \
           --log-opt max-file=1                                             \
           --restart=always                                                 \
           -p 127.0.0.1:8802:8801                                           \
           toomore/mailbox:cmd mailbox_server
