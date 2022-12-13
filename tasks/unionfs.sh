#!/bin/bash

backupDir=${1}
mountDir=${2}

mountBranches=$(find ${backupDir} -mindepth 1 -maxdepth 1 -type d -exec echo -n {}"=RO:" \; | sort -r)
mountBranches=${mountBranches:0:${#mountBranches}-1}

mountpoint -q ${mountDir}
if [ "$?" -eq 0 ]; then
    umount ${mountDir}
fi

unionfs ${mountBranches} ${mountDir}
