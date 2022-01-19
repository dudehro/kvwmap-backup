package create

import (
    "kvwmap-backup/docker"
	"kvwmap-backup/config"
    "kvwmap-backup/pgsql"
    "fmt"
    "os"
    "bufio"
    "strings"
    "strconv"
)

func readstdin(msg string) (string) {
    fmt.Printf(msg)
    var reader = bufio.NewReader(os.Stdin)
    ret,_ := reader.ReadString('\n')
    ret = strings.TrimSpace(ret)
//    fmt.Println(ret)
    return ret
}

func ask(msg string, fallback string) (string) {
    answer := readstdin(msg)
    if len(answer) == 0 {
        return fallback
    } else {
        return answer
    }
}

func New() {

    backup := config.NewConfig( readstdin("Speicherort für Backup: ") )

    networks := docker.ListNetworks()
    for _,network := range networks {
        if ask( fmt.Sprintf("Netzwerk %s sichern? [J=default/Nn]: ", network.Name), "J" ) == "J" {

            network_details := docker.InspectNetwork(network.Name)

            subnet := ""
            if len(network_details.IPAM.Config) != 0 {
                subnet = network_details.IPAM.Config[0].Subnet
            }
             currentNetworkPtr := config.AddNetwork2Backup(backup, network.Name, subnet)

            for _, c := range network_details.Containers {
                docker_container := docker.InspectContainer(c.Name)

                if ask(fmt.Sprintf("Container %s sichern? [J=default/Nn]: ", docker_container.Name), "J" ) == "J" {

                    currentContainerPtr := config.AddContainer2Backup(backup, docker.InspectImage(docker_container.Image).RepoTags[0], docker_container.Name, subnet)
                    config.AddNetwork2Container(currentContainerPtr, currentNetworkPtr.Name)

                    for _, mount := range docker_container.Mounts {
                        if config.IsMountInBackup(backup, mount.Source) {
                            fmt.Sprintf("Mount %s wird bereits mit anderen Containern gesicherte und hier übersprungen!", mount.Source)
                        } else {
                            if ask(fmt.Sprintf("Mount %s:%s sichern? [J=default/Nn]: ", mount.Source, mount.Destination), "J") == "J" {
                                config.AddMount(backup, currentContainerPtr, mount.Source, mount.Destination)
                            }
                        }
                    }

                    if strings.Contains(currentContainerPtr.Image, "pkorduan/postgis") || strings.Contains(currentContainerPtr.Image, "postgres") {
                        createPostgres(backup, currentContainerPtr, network_details.Containers[ docker_container.ID ].IPv4Address )
                    }

                    if strings.Contains(currentContainerPtr.Image, "mysql") {
                        createMysql(backup, currentContainerPtr/*, network_details.Containers[ docker_container.ID ].IPv4Address*/ )
                    }

                }
            }


        }
    }

    config.WriteFile("backup.json", backup)

}

func createPostgres(backup *config.Backup, container *config.Service, host string) {
    port := ask("Postgresql Port [Default=5432]: ", "5432")
    user := ask("Postgresql User [Default=postgres]: ", "postgres")
    password := ask("Postgresql Passwort [Default='']: ","")   //.pgpass sollte eingerichtet sein
    dbname := ask("Postgresql Datenbank [Default=postgres]: ","postgres")

    pgsql.OpenConnection("localhost", port, user, password, dbname)

    pg_dump := config.AddPostgres2Backup(backup, container, dbname, user, host)

    for _, schema := range pgsql.ListSchemas() {
        if ask(fmt.Sprintf("Schema %s sichern? [J=Default/N]", schema), "J") == "J" {
            config.AddSchema2PgDump(pg_dump, schema)
        }
    }

    if ask("Soll ein pg_dumpall-Eintrag für Rollen+Tablespaces dieses Clusters eingerichtet werden? [Default=j]: ","J") == "J" {
        config.AddPgDumpAll2Backup(backup, container, dbname, user, host, []string{"--globals-only"})
    }
}

func createMysql(backup *config.Backup, container *config.Service/*, host string*/) {
    var mysql *config.Mysql
    if ask("Sollen die Datenbanken 'kvwmap' und 'mysql' gesichert werden? [Default=J]: ", "J") == "J" {
        mysql = config.AddMysql2Backup(backup, container, "kvwmap", "", []string{"kvwmap","mysql"}, []string{container.Name}, []string{})
    }

    if ask("Sollen weitere Datenbanken in diesem Container gesichert werden? [Default=N]: ", "N") == "J" {
        for i,_ := strconv.Atoi( ask("Wie viele Datenbanken wollen Sie erfassen?: ", "0") ); i > 0; i-- {
            dbName := ask("Geben Sie den Namen der Datenbank ein: ","")
            config.AddMysqlDB2Container(mysql, container, dbName)
        }
    }
}
