#!/usr/bin/env bash
docker run -it --rm                                                            \
           --link mailbox-mariadb-prod:MARIADB                                 \
           -v $(pwd)/csv:/csv                                                  \
           -e "mailbox_ses_api=???"                                            \
           -e "mailbox_ses_key=???"                                            \
           -e "mailbox_ses_sender=???"                                         \
           -e "mailbox_web_site=???"                                           \
           toomore/mailbox:cmd sh
