package docker

import (
    "context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func ListContainers() ([]types.Container) {
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

func ListNetworks() ([]types.NetworkResource) {
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

func InspectContainer(containerID string) (types.ContainerJSON) {
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

func InspectNetwork(networkID string) (types.NetworkResource) {
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

func InspectImage(imageID string) (types.ImageInspect) {
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
