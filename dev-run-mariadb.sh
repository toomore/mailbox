#!/usr/bin/env bash
docker run -d --name mailbox-mariadb -v /srv/mailbox_mariadb:/var/lib/mysql \
           --log-opt max-size=64m                                           \
           --log-opt max-file=1                                             \
           -v $(pwd)/mariadb.cnf:/etc/mysql/conf.d/mariadb.cnf              \
           -e MYSQL_ROOT_PASSWORD=mailboxdbs                                \
           -e CHARACTER_SET_SERVER='utf8'                                   \
           -e COLLATION_SERVER='utf8_general_ci'                            \
           -e INIT_CONNECT='SET NAMES utf8'                                 \
           mariadb:10.1.24
