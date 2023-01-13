#!/bin/bash

backupDir=${1}
mountDir=${2}

mountBranches=""
i=0
while read SNAPSHOT
do
    ((i=i+1))
    if [ -z "$mountBranches" ]; then
        mountBranches=${SNAPSHOT}
    else
        mountBranches=${mountBranches}:${SNAPSHOT}
    fi

    if [ -f ${SNAPSHOT}/full.backup ]; then
        break 1
    fi
done < <(find ${backupDir} -mindepth 1 -maxdepth 1 -type d|sort -r)

mountpoint -q ${mountDir}
if [ "$?" -eq 0 ]; then
    umount ${mountDir}
fi

if [[ "$i" -eq 1 ]]; then
    MOUNTCMD="mount --bind ${mountBranches} ${mountDir}"
else
    MOUNTCMD="mount -t overlay overlay -o lowerdir=${mountBranches} ${mountDir}"
fi

echo "$MOUNTCMD"
$MOUNTCMD
