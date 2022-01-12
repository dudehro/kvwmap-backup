package main

import (
    "kvwmap-backup/logging"
    "kvwmap-backup/docker"
    "kvwmap-backup/create"
    "flag"
    "fmt"
)

func main(){

    log := logging.NewFileLogger()
    if err := log.StartLog(); err != nil {
        panic(err.Error())
    }

    modeFlag := flag.String("mode", "test", "Mode to run the tool, can be create, backup, test")

    flag.Parse()

    if *modeFlag == "create" {
        fmt.Println("create-Mode")
        create.New()
    } else if *modeFlag == "backup" {
        fmt.Println("backup-Mode")
    } else if *modeFlag == "test" {
        networks := docker.ListNetworks()
        for _,network := range networks {
            fmt.Printf("Network %s\n", network.Name)
            for _, c := range docker.InspectNetwork(network.Name).Containers {
                container := docker.InspectContainer(c.Name)
                fmt.Printf("Container: %s\n", container.Name)
                for _, mount := range container.Mounts {
                    fmt.Printf("Mount: %s:%s\n", mount.Source, mount.Destination)
                }
           }
        }
    } else {
        fmt.Println("no valid mode")
    }
}
