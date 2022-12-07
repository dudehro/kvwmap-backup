#!/bin/bash

QUELLE=/home/georgk/rsyncttest/source/
ZIEL=/home/georgk/rsyncttest/dest/


ZIELPOSTFIX_HEUTE=$(date +%F)/verzeichnisse
ZIELPOSTFIX_GESTERN=$(date -d "yesterday 13:00" '+%F')/verzeichnisse

ZIEL_LAST=${ZIEL}${ZIELPOSTFIX_GESTERN}
ZIEL_NOW=${ZIEL}${ZIELPOSTFIX_HEUTE}

echo "Quelle: ${QUELLE}"
echo "Ziel gestern: ${ZIEL_LAST}"
echo "Ziel heute: ${ZIEL_NOW}"

rsync -avP --del --include-from= --link-dest=${ZIEL_LAST} --log-file=${ZIEL_NOW}/rsync.$(date +%F).log ${QUELLE} ${ZIEL_NOW}
