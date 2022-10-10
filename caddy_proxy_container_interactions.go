package main

import (
	"context"
	"strings"

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

// GetServicesFromDocker gets services to reverse-proxy specified through Docker container labels
func GetServicesFromDocker(ctx context.Context, cli *client.Client) []utility.ServiceInfo {
	labelKeyValue := []string{
		LabelKeyForServerContainers(),
		LabelValueForServiceContainer(),
	}
	containers := utility.GetContainersByLabel(ctx, cli, labelKeyValue)

	services := []utility.ServiceInfo{}
	for _, container := range containers {
		ports := []uint16{}
		for _, portInfo := range container.Ports {
			// Note: Using PrivatePort since service & caddy-proxy-organizer is on the same network
			//    - So PrivatePort is port where the service is listening on
			ports = append(ports, portInfo.PrivatePort)
		}

		addressInExternalNetwork := ""
		for k,v := range container.NetworkSettings.Networks {
			if strings.Contains(k, NetworkNameForCaddyProxyExternal()) {
				addressInExternalNetwork = v.IPAddress
			}
		}
				
		serviceInfo := utility.ServiceInfo{
			Name: container.Names[0][1:], // Note: Need [1:] since name starts with "/"
			Domain: container.Labels[LabelKeyForServiceDomain()],
			Address: addressInExternalNetwork,
			Ports: ports,
		}
		services = append(services, serviceInfo)
	}

	return services
}

// ParseCaddyConfigTemplate parse the caddy json config template into a json string
func ParseCaddyConfigTemplate(ctx context.Context, cli *client.Client) string {
	cti := utility.ConsulTemplateImitator{
		Services: GetServicesFromDocker(ctx, cli),
	}
	parsedTemplate := cti.ParseTemplate(pathCaddyProxyConfigJsonTemplate())
	
	return parsedTemplate
}

// Generate a caddy json config from the template, and load it onto the caddy server
func LoadCaddyProxyJson(ctx context.Context, cli *client.Client) {
	caddyContainer := GetCaddyProxyContainer(ctx, cli)
	execConfig := types.ExecConfig{
		AttachStdin:  true,
		AttachStderr: true,  
		Cmd: ApiCaddyLoadCmd("", ParseCaddyConfigTemplate(ctx, cli)),
	}

	err := utility.RunDockerExec(ctx, cli, caddyContainer.ID, execConfig)
	if err != nil {
		panic(err)
	}
}