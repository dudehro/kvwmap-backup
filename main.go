package main

import (
	"flag"
	"fmt"
	"kvwmap-backup/create"
	"kvwmap-backup/docker"
	"kvwmap-backup/logging"
)

func main() {

	log := logging.NewFileLogger()
	if err := log.StartLog(); err != nil {
		panic(err.Error())
	}

	modeFlag := flag.String("mode", "ls", "Mode to run the tool, can be create, backup, ls")

	flag.Parse()

	if *modeFlag == "create" {
		fmt.Println("create-Mode")
		create.New()
	} else if *modeFlag == "backup" {
		fmt.Println("backup-Mode")
	} else if *modeFlag == "ls" {
		networks := docker.ListNetworks()
		for _, network := range networks {
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
