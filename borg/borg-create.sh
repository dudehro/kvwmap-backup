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
borgcmd="borg create -n --list --one-file-system -v --stats --progress --exclude-caches "

if [ -n "$patternsfrom" ]; then
    borgcmd="$borgcmd --patterns-from=$patternsfrom "
else
    echo "patternsfrom fehlt! Abbruch."
    exit 1
fi

if [ -n "$repopath" ]; then
    borgcmd="$borgcmd $repopath::{hostname}-{now:%Y-%m-%d-%H%M%S} "
else
    echo "repopath fehlt! Abbruch."
    exit 1
fi

echo "$borgcmd"
#$borgcmd
