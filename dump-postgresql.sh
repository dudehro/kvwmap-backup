#!/bin/bash

source ${1}
export CONTAINER
export DBUSER

EXITCODE=0

[ ! -d ${TARGET} ] && mkdir -p ${TARGET}

for DB in "${DATABASES[@]}"
do
    docker exec ${CONTAINER} bash -c "pg_dump -Fc -U ${DBUSER} -f /var/lib/postgresql/data/${CONTAINER}.${DB}.dump ${DB}"
    if [ $? -ne 0 ]; then
        EXITCODE=1
    fi
done

docker exec $CONTAINER bash -c "pg_dumpall -U $DBUSER -l postgres --globals-only -f /var/lib/postgresql/data/${CONTAINER}.schema_rollen.dump"
if [ $? -ne 0 ]; then
    EXITCODE=2
fi

DUMPDIR=$(docker inspect --format "{{json .Mounts}}" ${CONTAINER} | jq -r '.[]|select(.Destination=="/var/lib/postgresql/data").Source')
mv ${DUMPDIR}/*.*.dump ${TARGET}
if [ $? -ne 0 ]; then
    EXITCODE=3
fi

exit ${EXITCODE}
