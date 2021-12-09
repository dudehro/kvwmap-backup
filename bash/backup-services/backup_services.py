#!/usr/bin/env python3

import docker
import json
import os
import tarfile
from pathlib import Path

NETWORKS = [ "docker_static-network" ]
BACKUP_PATH = "/home/work/Sicherungen/networks"

docker = docker.from_env()
nw_list = docker.networks.list( names=NETWORKS )

# 1. Netzwerke auflisten
for network in nw_list:

    network = docker.networks.get(network.name)
    print("Sichere Netzwerk: " + network.name)
    network_path = os.path.join(BACKUP_PATH, network.name)
    Path( network_path ).mkdir( parents=True, exist_ok=True )
    print(network.attrs, file=open( os.path.join(network_path, 'docker-network-inspect'), 'a'))

    # 2. Container iterieren
    for container in network.containers:
        print("Container: " + container.name)
        service_path = os.path.join(network_path, 'services', container.name)
        Path( service_path ).mkdir( parents=True, exist_ok=True )
        print(container.attrs, file=open( os.path.join(service_path, 'docker-inspect'), 'a'))

        # Container als Image sichern
        image = container.attrs['Image']
        image = docker.images.get(image)
#        print( image.attrs['RepoTags'] )
        containerimg = container_tar = container.export()
        with open( os.path.join(service_path, "container.img"), "wb") as file:
            for chunk in containerimg:
                file.write(chunk)

        # 3. Mounts sichern

        for mount in container.attrs['Mounts']:
            print("sichere Mount: " + mount['Destination'])
            tarname = mount['Destination'].replace('/','_')
            print(mount['Source'] + ":" + mount['Destination'], file=open( os.path.join(service_path, 'tars_container_mounts'), 'a'))
#            with tarfile.open( os.path.join(service_path, tarname + '.tar.gz'), mode='w:gz') as archive:
#                archive.add( mount['Source'] , recursive=True)

