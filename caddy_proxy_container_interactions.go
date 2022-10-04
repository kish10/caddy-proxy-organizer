package main

import (
	"bytes"
	"context"
	"text/template"

	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
	"github.com/kish10/caddy-proxy-organizer/utility"
)

func GetCaddyProxyContainer(ctx context.Context, cli *client.Client) types.Container {
	labelKeyValue := []string{
		LabelKeyForServerContainers(),
		LabelValueForCaddyProxyContainer(),
	}

	containers := utility.GetContainersByLabel(ctx, cli, labelKeyValue)
	if len(containers) == 0 {
		utility.ErrorLog.Fatalf("\nCouldn't find caddy-proxy container with labelkey `%s` & labelValue `%s`", labelKeyValue[0], labelKeyValue[1])
	}

	return containers[0]
}

func parseCaddyConfigTemplate() string {
	templateJson, err := template.ParseFiles(pathCaddyProxyConfigJsonTemplate())
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	templateJson.Execute(&buf, nil)

	return buf.String()
}

func LoadCaddyProxyJson(ctx context.Context, cli *client.Client) {
	caddyContainer := GetCaddyProxyContainer(ctx, cli)
	execConfig := types.ExecConfig{
		AttachStdin:  true,
		AttachStderr: true,  
		Cmd: ApiCaddyLoadCmd("", parseCaddyConfigTemplate()),
	}

	err := utility.RunDockerExec(ctx, cli, caddyContainer.ID, execConfig)
	if err != nil {
		panic(err)
	}
}