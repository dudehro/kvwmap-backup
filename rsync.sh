#!/bin/bash

source ${1}

SIGNAL_FULL="full.backup"
SIGNAL_INCR="incremental.backup"

RSYNC_ARGS=(-aqP --stats --log-file=${ZIEL}/rsync.$(date +%F).log ${QUELLE} ${ZIEL})

if [ -d ${COMPARE} ] && [ ${MODUS} == "inkr" ]; then
        RSYNC_ARGS+=(--compare-dest=${COMPARE})
        BACKUP_SIGNAL=${SIGNAL_INCR}
else
        BACKUP_SIGNAL=${SIGNAL_FULL}
fi

if [ -n "$FILTER" ]; then
        RSYNC_ARGS+=("--filter=. ${FILTER}")
fi

echo "rsync "${RSYNC_ARGS[@]}
rsync "${RSYNC_ARGS[@]}"

if [ -f ${ZIEL}/${SIGNAL_FULL} ]; then
        rm ${ZIEL}/${SIGNAL_FULL}
fi
if [ -f ${ZIEL}/${SIGNAL_INCR} ]; then
        rm ${ZIEL}/${SIGNAL_INCR}
fi

touch ${ZIEL}/${BACKUP_SIGNAL}
