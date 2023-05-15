#!/bin/bash

Config=/etc/backup/borg.conf
if [ -n "${1}" ]; then
    Config=${1}
fi

if [ ! -f "$Config" ]; then
    echo "Konfiguration $Config existiert nicht. Bitte Pfad zur Konfiguration übergeben."
    exit 1
fi

source "$Config"

if [ -z "$repopath" ]; then
    echo "repopath nicht definiert. Abbruch!"
    exit 1
fi

borg init --encryption=none "$repopath"
