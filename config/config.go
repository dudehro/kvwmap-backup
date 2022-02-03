package config

import (
	//    "fmt"
	"encoding/json"
	"io/ioutil"
	"kvwmap-backup/docker"
	"log"
	"os"
)

func GetConfig(file string) (*Backup, bool) {
	_, err := os.Stat(file)
	if len(file) == 0 || os.IsNotExist(err) {
		p := Backup{}
		return &p, true
	} else {
		var p Backup
		file, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatalf("Datei %s kann nicht gelesen werden. Fehler: %s", file, err)
		}
		err = json.Unmarshal([]byte(file), &p)
		if err != nil {
			log.Fatalf("Datei %s entspricht nicht dem Schema. Fehler: %s", file, err)
		}
		return &p, false
	}
}

func AddNetwork2Backup(backup *Backup, networkName string, subnet string) *Network {
	if backup == nil {
		log.Fatal()
	}
	var networkPtr *Network
	_, networkPtr = IsNetworkInBackup(backup, networkName)
	if networkPtr == nil {
		//        fmt.Printf("Netwerk nict %s gefunden\n", networkName)
		networkPtr = &Network{Name: networkName, Subnet: subnet}
		backup.Networks = append(backup.Networks, networkPtr)
	} else {
		//        fmt.Printf("Netwerk %s gefunden: %+v\n", networkName, networkPtr)
		networkPtr.Subnet = subnet
	}
	return networkPtr
}

func IsNetworkInBackup(backup *Backup, networkName string) (bool, *Network) {
	for _, network := range backup.Networks {
		if network != nil {
			if network.Name == networkName {
				return true, network
			}
		}
	}
	return false, nil
}

func AddNetwork2Container(container *Service, networkname string) {
	if container == nil {
		log.Fatal()
	}
	found := false
	for _, network := range container.Networks {
		found = network == networkname
	}
	if !found {
		container.Networks = append(container.Networks, networkname)
	}
}

func AddContainer2Backup(backup *Backup, image string, name string, networks []string) *Service {
	if backup == nil {
		log.Fatal()
	}
	var service *Service
	found := false
	for _, container := range backup.Services {
		if container.Name == name {
			found = true
			container.Image = image
			container.Networks = networks
			service = container
			break
		}
	}
	if !found {
		service = &Service{Image: image, Name: name, Networks: networks}
		backup.Services = append(backup.Services, service)
	}
	return service
}

func AddMount(backup *Backup, container *Service, source string, destination string) *Mount {
	if backup == nil || container == nil {
		log.Fatal()
	}
	isthere, m := IsMountInBackup(backup, source)
	if isthere {
		return m
	} else {
		mount := Mount{MountSource: source, MountDestination: destination, Service: container.Name}
		backup.Mounts = append(backup.Mounts, &mount)
		return &mount
	}
}

func AddContainerConfig2Backup(backup *Backup, container *Service) {
	for _, m := range docker.GetContainerConfigFiles(container.Name) {
		AddMount(backup, container, m, "")
	}
}

func AddPostgresCluster2Backup(backup *Backup, container *Service, dbName string, dbUser string, dbHost string) *PgDump {
	if container == nil {
		log.Fatal()
	}
	var pgPtr *PgDump
	found := false
	if len(backup.PgDumps) > 0 {
		for _, pgdump := range backup.PgDumps {
			if pgdump.DbHost == dbHost && pgdump.DbName == dbName && pgdump.DbUser == dbUser {
				found = true
				pgPtr = pgdump
				break
			}
		}
	}
	if !found {
		pg := PgDump{DbName: dbName, DbUser: dbUser, DbHost: dbHost, Services: []string{container.Name}}
		backup.PgDumps = append(backup.PgDumps, &pg)
		pgPtr = &pg
	}
	return pgPtr
}

func AddSchema2PgDump(pgdump *PgDump, schema string) {
	pgdump.Schemas = append(pgdump.Schemas, schema)
}

func AddPgDumpAll2Backup(backup *Backup, container *Service, dbName string, dbUser string, dbHost string, parameters []string) *PgDumpall {
	if backup == nil || container == nil {
		log.Fatal()
	}
	pgdumpall := PgDumpall{DbHost: dbHost, DbName: dbName, DbUser: dbUser, Parameters: parameters, Services: []string{container.Name}}
	backup.PgDumpalls = append(backup.PgDumpalls, &pgdumpall)
	return &pgdumpall
}

func AddMysql2Backup(backup *Backup, container *Service, dbUser string, dbPassword string, databases []string, services []string, parameters []string) *Mysql {
	if backup == nil || container == nil {
		log.Fatal()
	}
	mysql := Mysql{DbUser: dbUser, DbPassword: dbPassword, Databases: databases, Services: services, Parameters: parameters}
	backup.Mysqls = append(backup.Mysqls, &mysql)
	return &mysql
}

func AddMysqlDB2Container(mysql *Mysql, container *Service, database string) {
	if mysql == nil || container == nil {
		log.Fatal()
	}
	mysql.Databases = append(mysql.Databases, database)
}

func WriteFile(location string, backup *Backup) {
	file, _ := json.MarshalIndent(backup, "", " ")
	_ = ioutil.WriteFile(location, file, 0644)
}

func IsMountInBackup(backup *Backup, source string) (bool, *Mount) {
	for _, mount := range backup.Mounts {
		if mount.MountSource == source {
			return true, mount
		}
	}
	return false, nil
}
