#!/bin/bash

SOURCE=/home/gisadmin/networks
RW_DIR=/home/gisadmin/Sicherungen/Verzeichnisse/rw
SNAPSHOT_DIR=/home/gisadmin/Sicherungen/Verzeichnisse/snapshots
SNAPSHOT_POSTFIX=$(date +%F)
MOUNTPOINT=/home/gisadmin/Sicherungen/Verzeichnisse-gesamt

mountpoint -q ${MOUNTPOINT}
if [ "$?" -eq 0 ]; then
    umount ${MOUNTPOINT}
    if [ "$?" -gt 0 ]; then
        echo "Fehler: Mountpoint ${MOUNTPOINT} konnte nicht ausgehangen werden."
        exit 1
    fi
fi

# Daten von rw nach ${SNAPSHOT_DIR}/${SNAPSHOT_POSTFIX} kopieren
mkdir ${SNAPSHOT_DIR}/${SNAPSHOT_POSTFIX}
rsync -avzP ${RW_DIR}/ ${SNAPSHOT_DIR}/${SNAPSHOT_POSTFIX}
if [ "$?" -gt 0 ]; then
        echo "Fehler beim verschieben der Daten!"
        exit 1
fi
rm -rdf ${RW_DIR}/*

mountBranches=$(find ${SNAPSHOT_DIR} -mindepth 1 -maxdepth 1 -type d -exec echo -n {}"=RO:" \; | sort -r)
mountBranches=${mountBranches:0:${#mountBranches}-1}
mountBranches="${RW_DIR}=RW:"${mountBranches}
unionfs -o cow ${mountBranches} ${MOUNTPOINT}
