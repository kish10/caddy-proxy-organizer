package main

import (
	"context"

	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
	"github.com/kish10/caddy-proxy-organizer/utility"
)

func GetCaddyProxyContainer(ctx context.Context, cli *client.Client) types.Container {
	labelKeyValue := []string{
		LabelKeyForServerContainers(),
		LabelValueForCaddyProxyContainer(),
	}
	utility.InfoLog.Printf("LabelKeyForServerContainers: %s", LabelKeyForServerContainers())

	containers := utility.GetContainersByLabel(ctx, cli, labelKeyValue)
	if len(containers) == 0 {
		utility.ErrorLog.Fatalf("\nCouldn't find caddy-proxy container with labelkey `%s` & labelValue `%s`", labelKeyValue[0], labelKeyValue[1])
	}

	return containers[0]
}