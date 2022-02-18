package main

import (
	"flag"
	"fmt"
	"kvwmap-backup/backup"
	"kvwmap-backup/create"
	"kvwmap-backup/docker"
	//	"kvwmap-backup/logging"
)

type mode string

const (
	ModeCreate string = "create"
	ModeBackup string = "backup"
	ModeList   string = "ls"
)

func main() {

	modeFlag := flag.String("mode", "ls", "Mode to run the tool, can be [create|backup|ls] ")
	configFlag := flag.String("backupconfig", "backup.json", "file to read and write the config to/from")
    dirFileFlag := flag.String("dirconfig", "", "file with directory structure")
	logFlag := flag.String("loglevel", "", "Loglevel: [info|warning|error|debug]")

	flag.Parse()

	if *modeFlag == ModeCreate {
		create.New(*configFlag)
	} else if *modeFlag == ModeBackup {
		backup.StartBackup(*configFlag, *dirFileFlag, *logFlag)
	} else if *modeFlag == ModeList {
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
