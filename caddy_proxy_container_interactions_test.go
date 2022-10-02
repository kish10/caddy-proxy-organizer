package main

import (
	"context"
	"testing"

	"github.com/kish10/caddy-proxy-organizer/utility"
)

const composeFilePath = "docker-compose--for-test.yaml"
var ctx = context.Background()

func TestGetCaddyProxyContainer(t *testing.T) {
	utility.RunDockerComposeUp(
		ctx, 
		utility.RunDockerComposeParams{[]string{composeFilePath}, false},
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