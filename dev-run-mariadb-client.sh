#!/usr/bin/env bash
docker run -it --rm --link mailbox-mariadb:MARIADB               \
           -v $(pwd)/sql:/sql                                    \
           -v $(pwd)/mariadb.cnf:/etc/mysql/conf.d/mariadb.cnf   \
           mariadb:10.11.10 bash
