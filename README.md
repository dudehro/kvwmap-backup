# tool for creating backups of docker-compose services in networks #

## flags ##
* -mode [create|backup|ls]
create - create new backup-configuration
backup - do a backup based on a configuration
ls - list all networks, containers and mounts, just for debugging
* -backupconfig [backup-config.json]
backup configuration to read and write in --mode [create|backup]
* -log [info|warning|error|debug]
sets the loglevel

