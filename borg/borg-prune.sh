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

borgcmd="borg prune --list -v "

if [ -z "$repopath" ]; then
    echo "repopath nicht gesetzt. Abbruch!"
    exit 1
fi

if [ -n "$keeplastdays" ]; then
    borgcmd="$borgcmd --keep-within=${keeplastdays}d "
fi

if [ -n "$keepweekly" ]; then
    borgcmd="$borgcmd --keep-weekly=$keepweekly "
fi

if [ -n "$keepmonthly" ]; then
    borgcmd="$borgcmd --keep-monthly=$keepmonthly"
fi

borgcmd="$borgcmd $repopath"
echo $borgcmd
$borgcmd
borg compact --cleanup-commits --progress "$repopath"
