#!/bin/bash

LOGFILE=${1}

TAR_COUNT=$(cat $LOGFILE | jq '.jobs | length')
if (( $TAR_COUNT > 0 )); then
    for (( i=0; i < $TAR_COUNT; i++ )); do
        name=$(cat $LOGFILE | jq -r ".jobs[$i].name")
	startT=$(cat $LOGFILE |jq -r ".jobs[$i].startime")
	endT=$(cat $LOGFILE |jq -r ".jobs[$i].endtime")
	startTf=$(date -d @$startT)
	endTf=$(date -d @$endT)
	exitcode=$(cat $LOGFILE |jq -r ".jobs[$i].exitcode")
	duration=$(expr $endT - $startT )
	echo $name
	echo -e '\t'"started: "$startTf
	echo -e '\t'"ended: "$endTf
	echo -e '\t'"duration: "$duration
	echo -e '\t'"exitcode: "$exitcode

	if [ $exitcode -gt 0 ]; then
		echo -e '\t'"stderr: "$(cat $LOGFILE |jq -r ".jobs[$i].stderr")
	fi


    done
fi

