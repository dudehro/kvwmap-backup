package create

import (
	"bufio"
	"fmt"
	"kvwmap-backup/config"
	"kvwmap-backup/docker"
	"kvwmap-backup/pgsql"
	"log"
	"os"
	"strconv"
	"strings"
    "time"
)

func readstdin(msg string) string {
	fmt.Printf(msg)
	var reader = bufio.NewReader(os.Stdin)
	ret, _ := reader.ReadString('\n')
	ret = strings.TrimSpace(ret)
	//    fmt.Println(ret)
	return ret
}

func ask(msg string, fallback string) string {
	answer := readstdin(msg)
	if len(answer) == 0 {
		return fallback
	} else {
		return answer
	}
}

func New(file string) {


	backup , newConfig := config.GetConfig(file)
    if newConfig {
        fmt.Printf("neue Backup-Config wird in Datei %s geschrieben\n", file)
    	backup.BackupPath = ask("statische Pfad für Sicherungen [/home/gisadmin/Sicherungen/woechentlich/ = Default]: ", "/home/gisadmin/Sicherungen/woechentlich")
//	    backup.BackupFolder = ask("Ordner für Backup, wird an statischen Pfad angehängt [2006_01_02 = Default]: ", "2006_01_02")
        backup.BackupFolder = fmt.Sprintf("%d-%02d-%02d", time.Now().Year(), time.Now().Month(), time.Now().Day())
  } else {
        fmt.Printf("bestehende Backup-Config aus Datei %s wird geändert\n", file)
    }

	networks := docker.ListNetworks()

	network_choice := "a"

	for network_choice != "" {
		fmt.Println("Liste der Netzwerke:")
		for n, network := range networks {
			fmt.Printf("%d) %s\n", n, network.Name)
		}
		network_choice = ask("Auswahl [ <leer> = Ende]: ", "")

		if network_choice == "" {
			break
		}

		choice_int, ok := strconv.Atoi(network_choice)
		if ok != nil {
			log.Fatal("keine gültige Zahl")
			break
		}

		network_details := docker.InspectNetwork(networks[choice_int].Name)

		subnet := ""
		if len(network_details.IPAM.Config) != 0 {
			subnet = network_details.IPAM.Config[0].Subnet
		}
		currentNetworkPtr := config.AddNetwork2Backup(backup, network_details.Name, subnet)
//        fmt.Printf("currentNetworkPtr %+v\n", currentNetworkPtr)

		for _, c := range network_details.Containers {
			docker_container := docker.InspectContainer(c.Name)

			if ask(fmt.Sprintf("Container %s sichern? [J=Default/N]: ", docker_container.Name), "J") == "J" {

				currentContainerPtr := config.AddContainer2Backup(backup, docker.InspectImage(docker_container.Image).RepoTags[0], docker_container.Name, []string{subnet})
				config.AddNetwork2Container(currentContainerPtr, currentNetworkPtr.Name)

				for _, mount := range docker_container.Mounts {
					if config.IsMountInBackup(backup, mount.Source) {
						fmt.Sprintf("Mount %s wird bereits mit anderen Containern gesicherte und hier übersprungen!", mount.Source)
					} else {
						if ask(fmt.Sprintf("Mount %s:%s sichern? [J=Default/Nn]: ", mount.Source, mount.Destination), "J") == "J" {
                            config.AddMount(backup, currentContainerPtr, mount.Source, mount.Destination)
                        }
					}
				}

				if strings.Contains(currentContainerPtr.Image, "pkorduan/postgis") || strings.Contains(currentContainerPtr.Image, "postgres") {
					createPostgres(backup, currentContainerPtr, network_details.Containers[docker_container.ID].IPv4Address)
				}

				if strings.Contains(currentContainerPtr.Image, "mysql") {
					createMysql(backup, currentContainerPtr /*, network_details.Containers[ docker_container.ID ].IPv4Address*/)
				}

			} //container j/n
		} //range containers
	} //network loop

	config.WriteFile(file, backup)

}

func createPostgres(backup *config.Backup, container *config.Service, host string) {
	port := ask("Postgresql Port [Default=5432]: ", "5432")
	user := ask("Postgresql User [Default=kvwmap]: ", "kvwmap")
	password := ask("Postgresql Passwort [Default='']: ", "") //.pgpass sollte eingerichtet sein
	dbname := ask("Postgresql Datenbank [Default=kvwmapsp]: ", "kvwmapsp")

	pgsql.OpenConnection("localhost", port, user, password, dbname)

	pg_dump := config.AddPostgresCluster2Backup(backup, container, dbname, user, host)

	for _, schema := range pgsql.ListSchemas() {
		if ask(fmt.Sprintf("Schema %s sichern? [J=Default/N]: ", schema), "J") == "J" {
			config.AddSchema2PgDump(pg_dump, schema)
		}
	}

	if ask("Soll ein pg_dumpall-Eintrag für Rollen+Tablespaces dieses Clusters eingerichtet werden? [J=Default/N]: ", "J") == "J" {
		config.AddPgDumpAll2Backup(backup, container, dbname, user, host, []string{"--globals-only"})
	}
}

func createMysql(backup *config.Backup, container *config.Service /*, host string*/) {
	var mysql *config.Mysql
	if ask("Sollen die Datenbanken 'kvwmap' und 'mysql' gesichert werden? [J=Default/N]: ", "J") == "J" {
		mysql = config.AddMysql2Backup(backup, container, "kvwmap", "", []string{"kvwmap", "mysql"}, []string{container.Name}, []string{})
	}

	if ask("Sollen weitere Datenbanken in diesem Container gesichert werden? [J/N=Default]: ", "N") == "J" {
		for i, _ := strconv.Atoi(ask("Wie viele Datenbanken wollen Sie erfassen?: ", "0")); i > 0; i-- {
			dbName := ask("Geben Sie den Namen der Datenbank ein: ", "")
			config.AddMysqlDB2Container(mysql, container, dbName)
		}
	}
}
