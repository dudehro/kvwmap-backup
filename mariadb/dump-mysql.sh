#!/bin/bash

source ${1}
export CONTAINER

MYSQLUSER=$(grep MYSQL_USER "$CREDENTIALS_FILE" | cut -d "'" -f 4)
MYSQLPW=$(grep MYSQL_PASSWORD "$CREDENTIALS_FILE" | cut -d "'" -f 4)
EXITCODE=0

[ ! -d ${TARGET} ] && mkdir -p ${TARGET}

for DB in "${DATABASES[@]}"
do
    FILENAME=${CONTAINER}.${DB}.dump
    docker exec "${CONTAINER}" bash -c "mysqldump -h mysql --single-transaction --user=$MYSQLUSER --databases ${DB} --password=\"$MYSQLPW\" > /var/lib/mysql/${FILENAME}"
    if [ $? -ne 0 ]; then
        EXITCODE=1
    fi
done

DUMPDIR=$(docker inspect --format "{{json .Mounts}}" ${CONTAINER} | jq -r '.[]|select(.Destination=="/var/lib/mysql").Source')
mv ${DUMPDIR}/${CONTAINER}.*.dump ${TARGET}
if [ $? -ne 0 ]; then
    EXITCODE=3
fi

exit ${EXITCODE}
