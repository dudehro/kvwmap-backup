#!/bin/bash

CONTAINER=${1}
BACKUPTYP=${2}

EXITCODE=0

if [ -z $CONTAINER ] || [ "$BACKUPTYP" != "full" -a "$BACKUPTYP" != "diff" -a "$BACKUPTYP" != "incr" ]; then
	echo "Aufruf: ./pgbackrest.sh [container] [full|diff|incr]"
	exit 1
fi

docker exec -it ${CONTAINER} pgbackrest --stanza=local --log-level-console=info --type=${BACKUPTYP} backup
if [ $? -ne 0 ]; then
    EXITCODE=2
fi

exit ${EXITCODE}
