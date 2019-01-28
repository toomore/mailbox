#!/usr/bin/env bash
docker run -it --rm --link mailbox-mariadb:MARIADB                             \
           -v $(pwd)/workdir:/workdir                                          \
           -v $(pwd)/campaign:/go/src/github.com/toomore/mailbox/campaign      \
           -v $(pwd)/cmd:/go/src/github.com/toomore/mailbox/cmd                \
           -v $(pwd)/mails:/go/src/github.com/toomore/mailbox/mails            \
           -v $(pwd)/reader:/go/src/github.com/toomore/mailbox/reader          \
           -v $(pwd)/utils:/go/src/github.com/toomore/mailbox/utils            \
           -v $(pwd)/main.go:/go/src/github.com/toomore/mailbox/main.go        \
           -p 127.0.0.1:8803:8801                                              \
           -e "mailbox_ses_key=???"                                            \
           -e "mailbox_ses_token=???"                                          \
           -e "mailbox_ses_sender=???"                                         \
           -e "mailbox_ses_replyto=???"                                        \
           -e "mailbox_web_site=???"                                           \
           golang:1.11.5 bash
