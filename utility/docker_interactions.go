package utility

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)


// GetDockerClient gets the Docker client
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

// ErrorExecCmdBadExit throws an error when the command run through Docker Exec exits ungracefully
type ErrorExecCmdBadExit struct {
	err error
	ContainerId string
	Cmd string
	CmdOutput string
	ExitCode int
}

func (e *ErrorExecCmdBadExit) Error() string {
    return fmt.Sprintf(
		"\nTried to run command:\n %s\nOn containerId:%s\nBut got exit code %d\nOutput from running the command:\n%s\n",
		e.Cmd,
		e.ContainerId,
		e.ExitCode,
		e.CmdOutput,
	)
}


// DockerExec executes a command on the given container
// Reference: 
// - https://github.com/moby/moby/blob/master/client/container_exec.go
// - https://stackoverflow.com/a/57132902
func RunDockerExec(ctx context.Context, cli *client.Client, containerId string, execConfig types.ExecConfig) error {
	if cli == nil {
		cli = GetDockerClient()
	}
	
	// NOTE (2022-10-02): At the moment these are not needed for ContainerExecAttach
	// DO NOT DELETE: Left as reference
	// DELETE WHEN : Discover that the settings are independent to ContainerExecAttach (or when comfortable) 
	// execConfig.AttachStdin = true
	// execConfig.AttachStderr = true


	// -- Run command on the container using Docker Exec
	InfoLog.Printf(
		"\nAbout to create Exec instance on container %s, with command:\n%s\n", 
		containerId,
		strings.Join(execConfig.Cmd, "\n"),
	)
	responseExecCreate, err := cli.ContainerExecCreate(ctx, containerId, execConfig)
	if err != nil {
		ErrorLog.Fatal(err)
	}

	responseExecAttach, errExecAttach := cli.ContainerExecAttach(ctx, responseExecCreate.ID, types.ExecStartCheck{}) 
	defer responseExecAttach.Close()
	if errExecAttach != nil {
		ErrorLog.Fatal(err)
	}
	scanner := bufio.NewScanner(responseExecAttach.Reader)
	attachBufferOutput := ""
	for scanner.Scan() {
		attachBufferOutput = attachBufferOutput + scanner.Text()
	}
	InfoLog.Print("\nExec Attach output:\n", attachBufferOutput)


	// Note (2022-10-02):  ContainerExecAttach, seems to also start the contaner
	// Reference: https://github.com/moby/moby/blob/master/client/container_exec.go
	// DELETE WHEN: Comfortable with only using ContainerExecAttach
	// err = cli.ContainerExecStart(ctx, responseExecCreate.ID, types.ExecStartCheck{})
	// if err != nil {
	// 	ErrorLog.Fatal(err)
	// }
	

	// -- Check whether Exec call exited gracefully 

	responseExecInspect, errExecInspect := cli.ContainerExecInspect(ctx, responseExecCreate.ID) 
	if errExecInspect != nil {
		ErrorLog.Fatal(errExecInspect)
	}
	
	if execExitCode := responseExecInspect.ExitCode; execExitCode > 0 {
		return &ErrorExecCmdBadExit{
			ContainerId: containerId,
			Cmd: strings.Join(execConfig.Cmd, " "),
			CmdOutput: attachBufferOutput,
			ExitCode: execExitCode,
		}
	} else {
		InfoLog.Printf("\nCommand ran via Docker Exec exited gracefully (exit code: %d) ", responseExecInspect.ExitCode)
	}

	return nil
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