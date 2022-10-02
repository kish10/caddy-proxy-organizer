package utility

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
func GetContainersByLabel(ctx context.Context, cli *client.Client, labelKeyValue []string) []types.Container {
	labelKey := labelKeyValue[0]
	labelValue := labelKeyValue[1]

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

type RunDockerComposeParams struct {
	ComposeFilePaths []string
	Rebuild bool
}

// RunDockerCompose up calls the shell process "docker compose -f <file_path> up -d"
func RunDockerComposeUp(ctx context.Context, args RunDockerComposeParams) {
	composefilePathFlags := []string{}
	for _, filePath := range args.ComposeFilePaths {
		composefilePathFlags = append(composefilePathFlags, "-f", filePath)
	}

	cmdArgs := []string{"docker", "compose"}
	cmdArgs = append(cmdArgs, composefilePathFlags...)
	cmdArgs = append(cmdArgs, "up", "-d")
	if args.Rebuild {
		cmdArgs = append(cmdArgs, "--build")
	}

	cmd := exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)
	InfoLog.Printf("\nAbout to run `docker compose up` with command:\n %s\n", cmd)
	stdoutStderr, err := cmd.CombinedOutput()
	if (err != nil) && (err.Error() != "exec: already started") {
		ErrorLog.Fatal(err)
	}
	InfoLog.Printf("\nResulting output from running `docker compose up` command:\n%s\n", stdoutStderr)
}


func StopContainersByLabel(ctx context.Context, cli *client.Client, labelKeyValue []string) {
	
	if cli == nil {
		cli = GetDockerClient()
	}

	for _, container := range GetContainersByLabel(ctx, cli, labelKeyValue) {
		if err := cli.ContainerStop(ctx, container.ID, nil); err != nil {
			panic(err)
		}
	}
}