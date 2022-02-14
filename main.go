package main

import (
	"flag"
	"fmt"
	"kvwmap-backup/backup"
	"kvwmap-backup/create"
	"kvwmap-backup/docker"
//	"kvwmap-backup/logging"
)

func main() {

	modeFlag := flag.String("mode", "ls", "Mode to run the tool, can be create, backup, ls")
	fileFlag := flag.String("file", "backup.json", "file to read and write the config to")
	logFlag := flag.String("loglevel", "", "Loglevel: info, warning, error, debug and combinations from those")

	flag.Parse()

	if *modeFlag == "create" {
		create.New(*fileFlag)
	} else if *modeFlag == "backup" {
		backup.StartBackup(*fileFlag, *logFlag)
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
