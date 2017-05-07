#!/usr/bin/env bash
docker run -it --rm  --link mailbox-mariadb:MARIADB              \
           -v $(pwd)/cmd:/cmd                                    \
           -v $(pwd)/campaign:/go/src/github.com/toomore/mailbox/campaign      \
           -v $(pwd)/utils:/go/src/github.com/toomore/mailbox/utils            \
           -p 127.0.0.1:8801:8801                                              \
           -e "mailbox_ses_api=???"                                            \
           -e "mailbox_ses_key=???"                                            \
           -e "mailbox_ses_sender=???"                                         \
           -e "mailbox_web_site=???"                                           \
           golang:1.8.1 bash
