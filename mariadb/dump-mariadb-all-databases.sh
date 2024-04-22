#!/bin/bash

source ${1}
export CONTAINER

MYSQLUSER=$(grep MYSQL_USER "$CREDENTIALS_FILE" | cut -d "'" -f 4)
MYSQLPW=$(grep MYSQL_PASSWORD "$CREDENTIALS_FILE" | cut -d "'" -f 4)
EXITCODE=0

[ ! -d ${TARGET} ] && mkdir -p ${TARGET}

FILENAME=${CONTAINER}.all-databases.dump
docker exec "${CONTAINER}" bash -c "mariadb-dump -h localhost --single-transaction --user=$MYSQLUSER --all-databases --password=\"$MYSQLPW\" > /var/lib/mysql/${FILENAME}"
if [ $? -ne 0 ]; then
    EXITCODE=1
fi

DUMPDIR=$(docker inspect --format "{{json .Mounts}}" ${CONTAINER} | jq -r '.[]|select(.Destination=="/var/lib/mysql").Source')
mv ${DUMPDIR}/${CONTAINER}.*.dump ${TARGET}
if [ $? -ne 0 ]; then
    EXITCODE=3
fi

exit ${EXITCODE}
