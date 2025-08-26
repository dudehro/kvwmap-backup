#!/bin/bash

docker run --rm -it \
--name postgres-list-backups \
-v $(pwd)/backup:/pgbackrest \
pkorduan/postgis:15-3.3 \
pgbackrest info
