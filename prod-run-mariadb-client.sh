#!/usr/bin/env bash
docker run -it --rm --link mailbox-mariadb-prod:MARIADB          \
           -v $(pwd)/sql:/sql                                    \
           -v $(pwd)/mariadb.cnf:/etc/mysql/conf.d/mariadb.cnf   \
           mariadb:10.1.22 bash
