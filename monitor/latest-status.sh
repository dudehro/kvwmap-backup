#!/bin/bash

#return codes
# 0 keine Fehler
# 1 läuft noch, nicht implementiert
# 2 Warnungen
# 3 Fehler

function set_exitcode() {
    if [ $1 -gt $max_exitcode ]; then
        max_exitcode=$1
    fi
}

Config=/etc/backup/jobs.json
if [ -n "${1}" ]; then
    Config=${2}
fi
if [ ! -f "$Config" ]; then
    echo "Konfiguration $Config existiert nicht. Bitte Pfad zur Konfiguration übergeben."
    exit 1
fi
workdir=$(dirname $(cat "$Config" | jq -r .workdir))
date_today=$(date +%F)
date_minus1day=$(date  --date="yesterday" +%F)
date_minus2days=$(date  --date="2 days ago" +%F)
lastWorkdir=
if [ -d "$workdir/$date_today" ]; then
    lastWorkdir="$workdir/$date_today"
elif [ -d "$workdir/$date_minus1day" ]; then
    lastWorkdir="$workdir/$date_minus1day"
elif [ -d "$workdir/$date_minus2days" ]; then
    lastWorkdir="$workdir/$date_minus2days"
else
    echo 3
    exit 3
fi

i_max=$(jq '.jobs|length' $lastWorkdir/joblog.json)
i=0
max_exitcode=0
while [ $i -lt $i_max ]
do
    exitcode=$(jq .jobs[$i].exitcode $lastWorkdir/joblog.json)
    jobname=$(jq -r .jobs[$i].name $lastWorkdir/joblog.json)
    if [ "$exitcode" -gt 0 ]; then
        if [[ $jobname =~ ^borg ]]; then
            if [ "$exitcode" -eq 1 ]; then
                set_exitcode 2
            else
                set_exitcode 3
            fi
        else
            set_exitcode 3
        fi
    fi

    ((i=i+1))
done
echo $max_exitcode
exit $max_exitcode

