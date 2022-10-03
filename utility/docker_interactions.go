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


type ErrorExecCmdBadExit struct {
	err error
	ContainerId string
	Cmd string
	CmdOutput string
	ExitCode int
}

func (e *ErrorExecCmdBadExit) Error() string {
    return fmt.Sprintf(
		"\nTried to run command:\n %s\n On containerId %s, but got exit code %d, run output:\n%s\n",
		e.ContainerId,
		e.Cmd,
		e.ExitCode,
		e.CmdOutput,
	)
}


// DockerExec executes a command on the given container
func RunDockerExec(ctx context.Context, cli *client.Client, containerId string, execConfig types.ExecConfig) error {
	if cli == nil {
		cli = GetDockerClient()
	}
	
	execConfig.AttachStdin = true
	execConfig.AttachStderr = true


	// -- Run Docker Exec
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
	// attachBufferOutput, _ := responseExecAttach.Reader.ReadString(io.EOF)
	InfoLog.Print("\nExec Attach output:\n", attachBufferOutput)

	err = cli.ContainerExecStart(ctx, responseExecCreate.ID, types.ExecStartCheck{})
	if err != nil {
		ErrorLog.Fatal(err)
	}
	

	// -- Check whether Exec call exited gracefully 

	responseExecInspect, errExecInspect := cli.ContainerExecInspect(ctx, responseExecCreate.ID) 
	if errExecInspect != nil {
		ErrorLog.Fatal(errExecInspect)
	}
	InfoLog.Printf("\nExec Exit Code: %d", responseExecInspect.ExitCode)
	if execExitCode := responseExecInspect.ExitCode; execExitCode > 0 {
		// ErrorLog.Fatalf("Executing command failed with exit code %d", execExitCode)
		return &ErrorExecCmdBadExit{
			ContainerId: containerId,
			Cmd: strings.Join(execConfig.Cmd, " "),
			CmdOutput: attachBufferOutput,
			ExitCode: execExitCode,
		}
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