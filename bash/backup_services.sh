#!/bin/bash

NETWORKS=("kvwmap_prod")
BACKUP_PATH=/home/gisadmin/Sicherungen/networks/

# 1. Netzwerke auflisten

while read NETWORK_NAME
do
	if [[ " ${NETWORKS[*]} " =~ " ${NETWORK_NAME} " ]]; then
		echo "Sichere Netzwerk $NETWORK_NAME"
		NETWORK_PATH=${BACKUP_PATH}/${NETWORK_NAME}
		mkdir -p ${NETWORK_PATH}/services/

		docker network inspect ${NETWORK_NAME} > ${NETWORK_PATH}/docker-inspect-network

		# 2. Services in den Netzwerken iterieren
		while read SERVICE_NAME
		do
			echo "Sichere Service $SERVICE_NAME"
			SERVICE_PATH=${NETWORK_PATH}/services/${SERVICE_NAME}
			mkdir -p ${SERVICE_PATH}

			docker inspect ${SERVICE_NAME} > ${SERVICE_PATH}/docker-inspect

			# 3. Daten sichern
			COUNT_MOUNT=$(docker inspect ${SERVICE_NAME} --format '{{json .Mounts}}' | jq '. | length')
			i=0
			while [ "$i" -lt "$COUNT_MOUNT" ];
			do
				# 3. Daten sichern
				MOUNT_SOURCE=$(docker inspect ${SERVICE_NAME} --format '{{json .Mounts}}' | jq -r ".[$i].Source")
				MOUNT_DESTINATION=$(docker inspect ${SERVICE_NAME} --format '{{json .Mounts}}' | jq -r ".[$i].Destination")
				FILENAME=$(echo ${MOUNT_DESTINATION} | tr '/' '_')".tar"
				echo "sichere ${MOUNT_SOURCE}"
				tar -cf ${SERVICE_PATH}/${FILENAME} ${MOUNT_SOURCE}
				i=$(($i+1))
			done < <(docker inspect ${SERVICE_NAME} --format '{{json .Mounts}}' | jq ".[$i]")

		done < <(docker network inspect ${NETWORK_NAME} --format '{{json .Containers}}' | jq -r '.[].Name')
	fi

done < <(docker network ls --format '{{.Name}}')
