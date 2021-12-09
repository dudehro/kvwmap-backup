#!/bin/bash

CONFIG_FILE=$1

#########################################################
## # JSON pruefen                                       #
#########################################################

dummy=$(cat $CONFIG_FILE | jq '.')
if [ $? -gt 0 ]; then
    ABORT_BACKUP=TRUE
    echo "Config-Datei ungültig!"
    exit 1
else
    ABORT_BACKUP=FALSE
fi

function pgdump(){
	echo "pgdump todo"
}

function mysqldump(){
	echo "mysql todo"
}

function targz(){
	echo "targz"
	COUNT_MOUNT=$(docker inspect ${SERVICE_NAME} --format '{{json .Mounts}}' | jq '. | length')
	i=0
	while [ "$i" -lt "$COUNT_MOUNT" ];
	do
		# 3. Daten sichern
		MOUNT_SOURCE=$(docker inspect ${SERVICE_NAME} --format '{{json .Mounts}}' | jq -r ".[$i].Source")
		MOUNT_DESTINATION=$(docker inspect ${SERVICE_NAME} --format '{{json .Mounts}}' | jq -r ".[$i].Destination")
		if [ "$(jq -r ".networks[] | select(.name==\"${NETWORK_NAME}") | .services[] | select(.name==\"${SERVICE_NAME}\") | .tar[] | select(.mount_destination==\"${MOUNT_DESTINATION}")" config.json | .save_data)" = true ]; then

			FILENAME=$(echo ${MOUNT_DESTINATION} | tr '/' '_')".tar"
			echo "sichere ${MOUNT_SOURCE}"
			tar -cf ${SERVICE_PATH}/${FILENAME} ${MOUNT_SOURCE}
			echo "${MOUNT_SOURCE}:${MOUNT_DESTINATION}:${FILENAME}" >> ${SERVICE_PATH}/tars_container_mounts

		fi
		i=$(($i+1))
	done < <(docker inspect ${SERVICE_NAME} --format '{{json .Mounts}}' | jq ".[$i]")
}

BACKUP_PATH=$(jq -r '.backup_path' ${CONFIG_FILE})
echo "Sicherung nach: ${BACKUP_PATH}"

# 1. Netzwerke iterieren
while read NETWORK_NAME
do
	if [ "$(jq -r ".networks[] | select(.name==\"${NETWORK_NAME}\") | .save_network" config.json)" = true ]; then

		echo "Sichere Netzwerk $NETWORK_NAME"
		NETWORK_PATH=${BACKUP_PATH}/${NETWORK_NAME}
		mkdir -p ${NETWORK_PATH}/services/

		docker network inspect ${NETWORK_NAME} > ${NETWORK_PATH}/docker-network-inspect

		# 2. Services in den Netzwerken iterieren
		while read SERVICE_NAME
		do
			# gibt es einen Eintrag für den Service in der Config?
			if [ "$(jq -r ".networks[] | select(.name==\"${NETWORK_NAME}\") | .services[] | select(.name==\"${SERVICE_NAME}\") | .save_service" config.json)" = true ]; then

				echo "Sichere Service $SERVICE_NAME"
				SERVICE_PATH=${NETWORK_PATH}/services/${SERVICE_NAME}
				mkdir -p ${SERVICE_PATH}

				docker inspect ${SERVICE_NAME} > ${SERVICE_PATH}/docker-inspect
				CONTAINER_IMAGE=$(docker inspect ${SERVICE_NAME} --format '{{.Config.Image}}' | cut -d ':' -f 1)

				pgdump
				mysqldump
				targz
			fi
		done < <(docker network inspect ${NETWORK_NAME} --format '{{json .Containers}}' | jq -r '.[].Name')
	fi

done < <(docker network ls --format '{{.Name}}')

echo "Fertig"
