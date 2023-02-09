#!/bin/bash

source ${1}
export CONTAINER

MYSQLUSER=$(grep MYSQL_USER "$CREDENTIALS_FILE" | cut -d "'" -f 4)
MYSQLPW=$(grep MYSQL_PASSWORD "$CREDENTIALS_FILE" | cut -d "'" -f 4)
EXITCODE=0

[ ! -d ${TARGET} ] && mkdir -p ${TARGET}

docker exec "${CONTAINER}" bash -c "mariabackup --backup --target-dir=/var/lib/mysql/backup --password=\"$MYSQLPW\" --user=$MYSQLUSER"
docker exec "${CONTAINER}" bash -c "mariabackup --prepare --target-dir=/var/lib/mysql/backup"
if [ $? -ne 0 ]; then
    EXITCODE=1
fi

DUMPDIR=$(docker inspect --format "{{json .Mounts}}" ${CONTAINER} | jq -r '.[]|select(.Destination=="/var/lib/mysql").Source')
mv ${DUMPDIR}/backup ${TARGET}
if [ $? -ne 0 ]; then
    EXITCODE=3
fi

exit ${EXITCODE}
