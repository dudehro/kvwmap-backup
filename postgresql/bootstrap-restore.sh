#!/bin/bash

WORKDIR=$(echo $WORKDIR | cut -d '/' -f 5-)
WORKDIR=/dumps/$WORKDIR/postgresql-dumps
docker exec -e WORKDIR=$WORKDIR kvwmap_prod_pgsql bash -c "/dumps/kvwmap-backup/postgresql/restore_backup_database.sh"
