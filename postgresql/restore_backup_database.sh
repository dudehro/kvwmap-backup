#!/bin/bash
printenv
echo "im container: $WORKDIR"
psql -h localhost -U pgadmin -d postgres -c "DROP DATABASE IF exists geobasis_backup ;"
psql -h localhost -U pgadmin -d postgres -c "CREATE DATABASE geobasis_backup WITH OWNER = gisadmin TEMPLATE=template0 ENCODING = 'UTF8' TABLESPACE = pg_default LC_COLLATE = 'en_US.utf8' LC_CTYPE = 'en_US.utf8' CONNECTION LIMIT = -1;"
pg_restore -h localhost -U pgadmin -d geobasis_backup $WORKDIR/kvwmap_prod_pgsql.geobasis_mse2.dump
psql -h localhost -U pgadmin -d geobasis_backup -f /dumps/config/spatial_ref_sys.sql
