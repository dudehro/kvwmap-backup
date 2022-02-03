package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"log"
	"os/exec"
    "strings"
    "path/filepath"
	//    "bufio"
)

func ListContainers() []types.Container {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	return containers
}

func ListNetworks() []types.NetworkResource {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	networks, err := cli.NetworkList(context.Background(), types.NetworkListOptions{})
	if err != nil {
		panic(err)
	}

	return networks
}

func InspectContainer(containerID string) types.ContainerJSON {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	container, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		panic(err)
	}
	return container
}

func GetContainerConfigFiles(containerID string) []string {
	var returnstr []string
	container := InspectContainer(containerID)
	working_dir := container.Config.Labels["com.docker.compose.project.working_dir"]
	config_files := strings.Split(container.Config.Labels["com.docker.compose.project.config_files"], ",")
	for _, s := range config_files {
		returnstr = append(returnstr, filepath.Join(working_dir, s))
	}
	return returnstr
}

func InspectNetwork(networkID string) types.NetworkResource {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	network, err := cli.NetworkInspect(context.Background(), networkID, types.NetworkInspectOptions{})
	if err != nil {
		panic(err)
	}

	return network
}

func InspectImage(imageID string) types.ImageInspect {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	image, _, err := cli.ImageInspectWithRaw(context.Background(), imageID)
	if err != nil {
		panic(err)
	}

	return image

}

func DockerExec(containerID string, cmd string, args ...string) string {
	args1 := []string{"exec", "-i", containerID, cmd}
	args2 := append(args1, args...)

	fmt.Println(exec.Command("docker", args2...).String())
	out, err := exec.Command("docker", args2...).Output()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(out))
	return string(out)
}

func Exec3(cmd string, args ...string) string {
	out, err := exec.Command(cmd, args...).Output()
	fmt.Println(string(out))
	if err != nil {
		log.Println(err)
	}
	return string(out)
}

/*
func containerExec(containerID string, cmd []string) error {
    cli, err := client.NewClientWithOpts(client.FromEnv)
    if err != nil {
        panic(err)
    }

	execOpts := types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
	}

	resp, err := cli.ContainerExecCreate(context.Background(), containerID, execOpts)
	if err != nil {
		return err
	}

	respTwo, err := cli.ContainerExecAttach(context.Background(), resp.ID, types.ExecConfig{})
	if err != nil {
		return err
	}
	defer respTwo.Close()

	err = cli.ContainerExecStart(context.Background(), resp.ID, types.ExecStartCheck{})
	if err != nil {
		return err
	}

	running := true
	for running {
		respThree, err := cli.ContainerExecInspect(context.Background(), resp.ID)
		if err != nil {
			panic(err)
		}

		if !respThree.Running {
			running = false
		}

        time.Sleep(250 * time.Millisecond)
	}

	return nil
}
*/
