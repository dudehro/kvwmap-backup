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
             _, currentNetworkPtr := config.AddNetwork(&backup, config.Network{Name: network.Name, SaveNetwork: true, Subnet: subnet} )

            for _, c := range network_details.Containers {
                docker_container := docker.InspectContainer(c.Name)

                if !config.IsContainerUnique(docker_container.Name, backup.Networks) {
                    break
                }

                if ask(fmt.Sprintf("Container %s sichern? [J=default/Nn]: ", docker_container.Name), "J" ) == "J" {

                    _, currentContainerPtr := config.AddContainer(currentNetworkPtr , config.Service{Name: docker_container.Name, Image: docker.InspectImage(docker_container.Image).RepoTags[0] })

                    for _, mount := range docker_container.Mounts {
                        if ask(fmt.Sprintf("Mount %s:%s sichern? [J=default/Nn]: ", mount.Source, mount.Destination), "J") == "J" {
							config.AddTarItem(currentContainerPtr, mount.Source, mount.Destination)
                        }
                    }

                    if strings.Contains(currentContainerPtr.Image, "pkorduan/postgis") || strings.Contains(currentContainerPtr.Image, "postgres") {
                        createPostgres(currentContainerPtr, network_details.Containers[ docker_container.ID ].IPv4Address )
                    }

                    if strings.Contains(currentContainerPtr.Image, "mysql") {
                        createMysql(currentContainerPtr/*, network_details.Containers[ docker_container.ID ].IPv4Address*/ )
                    }

                }
            }


        }
    }

    config.WriteFile("backup.json", backup)

}

func createPostgres(container *config.Service, host string) {
    port := ask("Postgresql Port [Default=5432]: ", "5432")
    user := ask("Postgresql User [Default=postgres]: ", "postgres")
    password := ask("Postgresql Passwort [Default='']: ","")   //.pgpass sollte eingerichtet sein
    dbname := ask("Postgresql Datenbank [Default=postgres]: ","postgres")

    pgsql.OpenConnection("localhost", port, user, password, dbname)

    pg_cluster := config.Postgres{DbName: dbname, DbUser: user, Host: host}
    container.Postgres = &pg_cluster
//    pgdump := config.Pgdump{IncludeListedSchemas: true}
    container.Postgres.Pgdump = &config.Pgdump{IncludeListedSchemas: true}

    for _, schema := range pgsql.ListSchemas() {
        if ask(fmt.Sprintf("Schema %s sichern? [J=Default/N]", schema), "J") == "J" {
            container.Postgres.Pgdump.Schemas = append(container.Postgres.Pgdump.Schemas, schema)
        }
    }

    pg_dumpall_item := config.PgdumpallItems{PgDumpallParameter: "--globals-only"}
    container.Postgres.Pgdumpall = append(container.Postgres.Pgdumpall, &pg_dumpall_item)

}

func createMysql(container *config.Service/*, host string*/) {
    if ask("Sollen die Datenbanken 'kvwmap' und 'mysql' gesichert werden? [Default=J]: ", "J") == "J" {
        container.Mysql = &config.Mysql{}
        container.Mysql.Mysqldump = append(container.Mysql.Mysqldump, &config.MysqldumpItems{DbName: "kvwmap"}, &config.MysqldumpItems{DbName: "mysql"} )
    }

    if ask("Sollen weitere Datenbanken in diesem Container gesichert werden? [Default=N]: ", "N") == "J" {
        for i,_ := strconv.Atoi( ask("Wie viele Datenbanken wollen Sie erfassen?: ", "0") ); i > 0; i-- {
            dbName := ask("Geben Sie den Namen der Datenbank ein: ","")
            container.Mysql.Mysqldump = append(container.Mysql.Mysqldump, &config.MysqldumpItems{DbName: dbName} )
        }
    }
}
