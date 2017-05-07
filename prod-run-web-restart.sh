#!/usr/bin/env bash
docker stop -t 3 mailbox-prod-web
docker rm mailbox-prod-web
sh ./prod-run-web.sh

docker stop -t 3 mailbox-prod-web-2
docker rm mailbox-prod-web-2
sh ./prod-run-web.sh
