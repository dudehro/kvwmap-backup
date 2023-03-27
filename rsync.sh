#!/bin/bash

# Scriptargumente + Config validieren
MODUS=${1}
if [ -z "$MODUS" ]; then
    echo "Argument fehlt: rsync.sh [full|diff|incr] [config]"
fi

Config=$HOME/.config/backup/backup.conf
if [ -n "${2}" ]; then
    Config=${2}
fi

if [ ! -f "$Config" ]; then
    echo "Konfiguration $Config existiert nicht. Bitte Pfad zur Konfiguration übergeben."
    exit 1
fi
# shellcheck source=/root/.config/backup/backup.conf
source "$Config"

if [ -z "$rsync_source" ] || [ -z "$rsync_dest" ] || [ -z "$OverlayDir" ]; then
    echo "Konfiguration unvollständig. Benötige rsync_source, rsync_dest, OverlayDir!"
    exit 1
fi
QUELLE=${rsync_source}
ZIEL=${rsync_dest}
FILTER=${rsync_filter}
COMPARE=${OverlayDir}

SIGNAL_FULL="full.backup"
SIGNAL_INCR="incremental.backup"
SIGNAL_DIFF="differential.backup"

# rsync Argumente

RSYNC_ARGS=("-aqP" "--stats" "--log-file=${ZIEL}/rsync.$(date +%F).log" "$QUELLE" "$ZIEL")

if [ -d ${COMPARE} ] && [ "$MODUS" == "incr" ]; then
    RSYNC_ARGS+=("--compare-dest=$COMPARE")
    BACKUP_SIGNAL=${SIGNAL_INCR}
elif [ -d ${COMPARE} ] && [ "$MODUS" == "diff" ]; then
    RSYNC_ARGS+=("--compare-dest=$COMPARE")
    BACKUP_SIGNAL=${SIGNAL_DIFF}
else
    BACKUP_SIGNAL=${SIGNAL_FULL}
fi

if [ -n "$FILTER" ]; then
    RSYNC_ARGS+=("--filter=. ${FILTER}")
fi

echo "rsync ${RSYNC_ARGS[@]}"
rsync "${RSYNC_ARGS[@]}"

# Liste mit Vorgängern/Abhängigkeiten nach $BACKUP_SIGNAL schreiben

if [ "$MODUS" == "incr" ]; then
    if [ -f "$COMPARE/$SIGNAL_INCR" ]; then
        VORGAENGERDATEI="$COMPARE/$SIGNAL_INCR"
    elif [ -f "$COMPARE/$SIGNAL_FULL" ]; then
        VORGAENGERDATEI="$COMPARE/$SIGNAL_FULL"
    else
        echo "Fehler: keinen Vorgänger bei inkrementeller Sicherung gefunden!"
        exit 1
    fi
elif [ "$MODUS" == "diff" ]; then
    if [ -f "$COMPARE/$SIGNAL_DIFF" ]; then
        VORGAENGERDATEI="$COMPARE/$SIGNAL_DIFF"
    elif [ -f "$COMPARE/$SIGNAL_FULL" ]; then
        VORGAENGERDATEI="$COMPARE/$SIGNAL_FULL"
    else
        echo "Fehler: keinen Vorgänger bei differenzieller Sicherung gefunden!"
        exit 1
    fi
elif [ "$MODUS" == "full" ]; then
    VORGANEGERDATEI=
fi

VORGAENGER=()
if [ -n "$VORGAENGERDATEI" ]; then
    readarray -t VORGAENGER < "$VORGAENGERDATEI"
fi

VORGAENGER+=("$ZIEL")
echo "${VORGAENGER[@]}" > "$ZIEL/$BACKUP_SIGNAL"
