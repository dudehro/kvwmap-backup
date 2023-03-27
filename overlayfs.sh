#!/bin/bash

MODUS=${1}
if [ -z "$MODUS" ]; then
    echo "Argument fehlt: overlayfs.sh [diff|incr] [config]"
fi

Config=$HOME/.config/backup/backup.conf
if [ -n "$2" ]; then
    Config="$2"
fi

if [ ! -f "$Config" ]; then
    echo "Konfiguration $Config existiert nicht."
    exit 1
fi

# shellcheck source=/root/.config/backup/backup.conf
source "$Config"

if [ -z "$SnapshotDir" ] || [ -z "$OverlayDir" ]; then
    echo "Konfiguration unvollständig. Benötige SnapshotDir und OverlayDir!"
    exit 1
fi

backupDir="$SnapshotDir"
mountDir="$OverlayDir"
DiffIncr="$MODUS"   #defaults to "incr"

if [ -z "$DiffIncr" ]; then
    DiffIncr="incr"
fi

# select snapshots
mountBranches=""
branchCount=0
while read SNAPSHOT
do
    if [ "$DiffIncr" = "diff" ]; then
        if [ -z "$mountBranches" ]; then
            ((branchCount=branchCount+1))
            mountBranches=${SNAPSHOT}
        fi

        if [ -f "$SNAPSHOT/full.backup" ]; then
            ((branchCount=branchCount+1))
            if [ -z "$mountBranches" ]; then
                mountBranches=${SNAPSHOT}
            else
                mountBranches=${mountBranches}:${SNAPSHOT}
            fi
        fi

    elif [ "$DiffIncr" = "incr" ]; then
        ((branchCount=branchCount+1))
        if [ -z "$mountBranches" ]; then
            mountBranches=${SNAPSHOT}
        else
            mountBranches=${mountBranches}:${SNAPSHOT}
        fi
    fi

    if [ -f "$SNAPSHOT/full.backup" ]; then
        break 1
    fi
done < <(find "$backupDir" -mindepth 1 -maxdepth 1 -type d|sort -r)

# unmount
if mountpoint -q "$mountDir"; then
    umount "$mountDir"
fi

# build command
if [[ "$branchCount" -eq 1 ]]; then
    MOUNTCMD="mount --bind ${mountBranches} ${mountDir}"
else
    MOUNTCMD="mount -t overlay overlay -o lowerdir=${mountBranches} ${mountDir}"
fi

echo "$MOUNTCMD"
$MOUNTCMD

