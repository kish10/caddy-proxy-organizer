package main

import (
	"context"
	"testing"

	"github.com/kish10/caddy-proxy-organizer/utility"
)

var ctx = context.Background()
const pathComposeFile = "./example_docker_compose/docker-compose--for-test.yaml"
// const pathCadyProxyJson = "caddy--for-test.json"

func TestGetCaddyProxyContainer(t *testing.T) {
	utility.RunDockerComposeUp(
		ctx, 
		utility.RunDockerComposeParams{[]string{pathComposeFile}, false},
	)

	caddyContainer := GetCaddyProxyContainer(ctx, nil)

	isCaddyContainer := false
	for k,v := range caddyContainer.Labels {
		if k == LabelKeyForServerContainers() && v == LabelValueForCaddyProxyContainer() {
			isCaddyContainer = true
		}
	}

	if !isCaddyContainer {
		t.Errorf(
			"No caddy-proxy container found with label key `%s` & value `%s`", 
			LabelKeyForServerContainers(),
			LabelValueForCaddyProxyContainer(),
		)
	}
}

func TestLoadCaddyProxyJson(t *testing.T) {
	// -- Test if compiles
	LoadCaddyProxyJson(ctx, nil)
}