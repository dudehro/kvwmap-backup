package create

import (
    "kvwmap-backup/docker"
	"kvwmap-backup/config"
    "fmt"
    "os"
    "bufio"
    "strings"
)

func readstdin(msg string) (string) {
    fmt.Printf(msg)
    var reader = bufio.NewReader(os.Stdin)
    ret,_ := reader.ReadString('\n')
    ret = strings.TrimSpace(ret)
//    fmt.Println(ret)
    return ret
}

func askJN(msg string, fallback string) (string) {
    answer := strings.ToUpper(readstdin(msg))
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
        if askJN( fmt.Sprintf("Netzwerk %s sichern? [Jj/Nn=default]: ", network.Name), "N" ) == "J" {

            _, currentNetworkPtr := config.AddNetwork(&backup, config.Network{Name: network.Name, SaveNetwork: true} )

            for _, c := range docker.InspectNetwork(network.Name).Containers {
                docker_container := docker.InspectContainer(c.Name)

                if !config.IsContainerUnique(docker_container.Name, backup.Networks) {
                    break
                }

                if askJN(fmt.Sprintf("Container %s sichern? [Jj=default/Nn]: ", docker_container.Name), "J" ) == "J" {

                    _, currentContainerPtr := config.AddContainer(currentNetworkPtr , config.Service{Name: docker_container.Name, Image: docker.InspectImage(docker_container.Image).RepoTags[0] })

                    for _, mount := range docker_container.Mounts {
                        if askJN(fmt.Sprintf("Mount %s:%s sichern? [Jj=default/Nn]: ", mount.Source, mount.Destination), "J") == "J" {
							config.AddTarItem(currentContainerPtr, mount.Source, mount.Destination)
                        }
                    }

                }
            }


        }
    }

    config.WriteFile("backup.json", backup)

}
