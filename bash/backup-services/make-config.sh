#!/bin/bash

CONFIG_FILE=$(pwd)/$1
echo $CONFIG_FILE
rm ${CONFIG_FILE}

function toconf(){
	echo ${1} >> ${CONFIG_FILE}
	echo ${1}
}

function tarconf(){

	toconf '"tar": {'
	toconf '  "diff_backup_days":"",'
	toconf '  "directories": [{}'

	COUNT_MOUNT=$(docker inspect ${SERVICE_NAME} --format '{{json .Mounts}}' | jq '. | length')
	i=0

	while [ "$i" -lt "$COUNT_MOUNT" ];
	do
		# 3. Daten sichern
		MOUNT_SOURCE=$(docker inspect ${SERVICE_NAME} --format '{{json .Mounts}}' | jq -r ".[$i].Source")
		MOUNT_DESTINATION=$(docker inspect ${SERVICE_NAME} --format '{{json .Mounts}}' | jq -r ".[$i].Destination")
		FILENAME=$(echo ${MOUNT_DESTINATION} | tr '/' '_')".tar"

		toconf '    ,{'
		toconf "      \"mount_source\":\"${MOUNT_SOURCE}\","
		toconf "      \"mount_destination\":\"${MOUNT_DESTINATION}\","
		toconf '      "save_data":true,'
		toconf '      "exclude_dirs":""'
		toconf '    }'

		i=$(($i+1))
	done < <(docker inspect ${SERVICE_NAME} --format '{{json .Mounts}}' | jq ".[$i]")

	toconf '  ]'
	toconf '}'
}

toconf '{'
toconf '  "backup_path":"/home/gisadmin/Sicherungen/networks",'
toconf '  "networks": [{}'

# 1. Netzwerke iterieren
while read NETWORK_NAME
do

	toconf ',{'
	toconf "  \"name\":\"${NETWORK_NAME}\","
        toconf '  "save_network":true,'
        toconf '  "services": [{}'

	# 2. Services in den Netzwerken iterieren
	while read SERVICE_NAME
	do
		toconf '    ,{'
		toconf "      \"name\":\"${SERVICE_NAME}\","
		toconf '      "save_service":true,'

		SERVICE_IMAGE=$(docker inspect ${SERVICE_NAME} --format '{{.Config.Image}}' | cut -d ':' -f 1)

		case $SERVICE_IMAGE in
			pkorduan/postgis)
				toconf '"postgres": {'
				toconf '  "db_user":"kvwmap",'
				toconf '  "db_name":"kvwmapsp",'
				toconf '  "pgdumpall": ['
				toconf '     {"pg_dumpall_parameter":""}'
				toconf '   ]'
				toconf '},'
			;;
			mysql)
				toconf '"mysql":'
				toconf '  {'
				toconf '    "db_user":"kvwmap",'
				toconf '    "db_password":"",'
				toconf '    "mysqldump" : ['
				toconf '       {'
				toconf '         "db_name":"mysql",'
				toconf '         "mysql_dump_parameter":""'
				toconf '       }'
				toconf '    ]'
				toconf '  },'
			;;
		esac

		tarconf

		toconf '    }'	#service-Objekt schließen
	done < <(docker network inspect ${NETWORK_NAME} --format '{{json .Containers}}' | jq -r '.[].Name')

	toconf '  ]'
	toconf '}'
done < <(docker network ls --format '{{.Name}}')

toconf '  ]'
toconf '}'

jq '.' ${CONFIG_FILE} > ${CONFIG_FILE}_tmp
mv ${CONFIG_FILE}_tmp ${CONFIG_FILE}
