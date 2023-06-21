#!/bin/bash
Config=/etc/backup/jobs.json
if [ -n "${1}" ]; then
    Config=${1}
fi

if [ ! -f "$Config" ]; then
    echo "Konfiguration $Config existiert nicht. Bitte Pfad zur Konfiguration übergeben."
    exit 1
fi

WORKDIR=$(dirname $(jq -r .workdir "$Config"))
find "$WORKDIR" -mindepth 1 -maxdepth 1 -type d -mtime +10 -exec rm -r {} \;
