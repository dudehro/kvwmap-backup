#!/bin/bash

backupDir=${1}
mountDir=${2}

mountBranches=$(find ${backupDir} -mindepth 1 -maxdepth 1 -type d -exec echo {} \;|sort -r|tr "\n" ":")
mountBranches=${mountBranches:0:${#mountBranches}-1}

mountpoint -q ${mountDir}
if [ "$?" -eq 0 ]; then
    umount ${mountDir}
fi

echo "mount -t overlay overlay -o lowerdir=${mountBranches} ${mountDir}"
mount -t overlay overlay -o lowerdir=${mountBranches} ${mountDir}
