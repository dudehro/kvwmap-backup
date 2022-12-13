#!/bin/bash

QUELLE=${1}
ZIEL=${2}
MODUS=${3}
FILTER=${4}

ZIELPOSTFIX_HEUTE=/$(date +%F)
ZIELPOSTFIX_GESTERN=/$(date -d "yesterday 13:00" '+%F')

ZIEL_LAST=${ZIEL}${ZIELPOSTFIX_GESTERN}
ZIEL_NOW=${ZIEL}${ZIELPOSTFIX_HEUTE}

RSYNC_ARGS=(-avP --del --log-file=${ZIEL_NOW}/rsync.$(date +%F).log ${QUELLE} ${ZIEL_NOW})

if [ -d ${ZIEL_LAST} ] && [ ${MODUS} == "inkr" ]; then
        RSYNC_ARGS+=(--compare-dest=${ZIEL_LAST})
fi

if [ -n "$FILTER" ]; then
        RSYNC_ARGS+=("--filter=. ${FILTER}")
fi

echo "rsync "${RSYNC_ARGS[@]}
rsync "${RSYNC_ARGS[@]}"

