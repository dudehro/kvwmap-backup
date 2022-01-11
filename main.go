package main

import (
    "kvwmap-backup/logging"
    "kvwmap-backup/docker"
    "fmt"
)

func main(){

    log := logging.NewFileLogger()
    if err := log.StartLog(); err != nil {
        panic(err.Error())
    }

    networks := docker.ListNetworks()
    for _,network := range networks {
        fmt.Printf("Network %s\n", network.Name)
        for _, c := range docker.InspectNetwork(network.Name).Containers {
            //fmt.Printf("Container: %s, %s\n", container.EndpointID[:10], container.Name)
            container := docker.InspectContainer(c.Name)
            fmt.Printf("Container: %s\n", container.Name)
            for _, mount := range container.Mounts {
                fmt.Printf("Mount: %s:%s\n", mount.Source, mount.Destination)
            }
       }
    }

}
