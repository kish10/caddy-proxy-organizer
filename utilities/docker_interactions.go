package utilities

import (
	"context"
	"os/exec"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// gets Docker client
func GetDockerClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	return cli
}

// GetContainersAll returns a list of running containers
func GetContainersAll(ctx context.Context, cli *client.Client) []types.Container {

	if cli == nil {
		cli = GetDockerClient()
	}
	
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	runningContainers := []types.Container{}
	for _, container := range containers {
		if container.State == "running" {
			runningContainers = append(runningContainers, container)
		}
	}

	return runningContainers
}

// GetContainersByLabel gets list of running containers with the given label identifiers
func GetContainersByLabel(ctx context.Context, cli *client.Client, labelArgs...string) []types.Container {
	labelKey := labelArgs[0]
	labelValue := labelArgs[1]

	needKeyAndValue := (labelKey != "") && (labelValue != "")

	containersWithLabel := []types.Container{}

	for _, container := range GetContainersAll(ctx, cli) {
		for k,v := range container.Labels {
			a := k == labelKey && !needKeyAndValue
			b := v == labelValue && !needKeyAndValue
			c := k == labelKey && v == labelValue
			switch {
			case a,b,c:
				containersWithLabel = append(containersWithLabel, container)
			}
		}
	}

	return containersWithLabel
}

// runDockerCompose up calls the shell process "docker compose -f <file_path> up -d"
func runDockerComposeUp(filePath string) {
	cmd := exec.Command("docker", "compose", "-f", filePath, "up", "-d")
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}


func StopContainersByLabel(ctx context.Context, cli *client.Client, labelArgs...string) {
	
	if cli == nil {
		cli = GetDockerClient()
	}

	for _, container := range GetContainersByLabel(ctx, cli, labelArgs...) {
		if err := cli.ContainerStop(ctx, container.ID, nil); err != nil {
			panic(err)
		}
	}
}