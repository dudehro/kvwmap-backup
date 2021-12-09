#!/bin/bash

#PGHOST
#PGUSER
#PGDATABASE

docker run \
-e PGHOST=kvwmap_prod_pgsql_1 \
-e PGUSER=kvwmap \
-e PGDATABASE=kvwmapsp \
-v /home/gisadmin/Sicherungen/scripte/kvwmap-backup/bash/dumps:/dumps \
--network="kvwmap_prod" \
--rm -it \
--name pg_dump \
gkaemmert/db_tools:latest \
/scripte/pgdump.sh
