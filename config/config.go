package config

import (
    "fmt"
    "log"
    "encoding/json"
    "io/ioutil"
)

func NewConfig(Path string) Backup {
    p := Backup{BackupPath: Path}
    return p
}

func AddNetwork(backup *Backup, network Network) (*Backup, *Network){
    if backup == nil {
        log.Fatal()
    }
    backup.Networks = append(backup.Networks, &network)
    return backup, &network
}

func AddContainer(network *Network, container Service) (*Network, *Service) {
    if network == nil {
        log.Fatal()
    }
    container.Tar = &Tar{}
    network.Services = append(network.Services, &container)
    return network, &container
}

func AddTarItem(container *Service, source string, destination string) (*Service) {
    if container == nil {
        log.Fatal()
    }
    container.Tar.Directories = append(container.Tar.Directories, &Taritem{MountSource: source, MountDestination: destination} )
    return container
}

func AddPostgresContainer(container *Service, dbName string, dbUser string) {
    if container == nil {
        log.Fatal()
    }
    pg := Postgres{DbName: dbName, DbUser: dbUser}
    container.Postgres = &pg
}

func IsContainerUnique(containerID string, networks []*Network) (bool) {
    containercount := 0
//    fmt.Printf("Type: %T, Values: %+v \n", networks, networks)

    for _,network := range networks {
        for _,service := range network.Services {
            if service.Name == containerID {
                containercount++
            }
        }
    }
    return containercount <= 1
}

func WriteFile(location string, backup Backup) {
    file, _ := json.MarshalIndent(backup, "", " ")
	_ = ioutil.WriteFile(location, file, 0644)
}

func Print(f string){
    fmt.Println(f)
}
